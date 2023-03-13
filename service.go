package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/helper/id"
	"github.com/UritMedical/qf/util/io"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

type Service struct {
	folder      string                    // 框架的文件夹路径
	db          *gorm.DB                  // 数据库
	engine      *gin.Engine               // gin
	bllList     map[string]IBll           // 所有创建的业务层对象
	apiHandler  map[string]ApiHandler     // 所有注册的业务API函数指针
	msgHandler  map[string]MessageHandler // 所有消息执行的函数指针
	setting     setting                   // 框架配置
	idAllocator iIdAllocator              // id分配器
}

//
// newService
//  @Description: 创建框架服务
//  @return *Service 服务对象指针
//
func newService() *Service {
	s := &Service{
		bllList:    map[string]IBll{},
		apiHandler: map[string]ApiHandler{},
		msgHandler: map[string]MessageHandler{},
		setting:    setting{},
	}
	// 默认文件夹路径
	s.folder = "."
	// 加载配置
	s.setting.Load(fmt.Sprintf("%s/config.toml", s.folder))
	// 创建数据库
	dbDir := io.CreateDirectory(fmt.Sprintf("%s/db", s.folder))
	gc := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	}
	if s.setting.GormConfig.OpenLog == 1 {
		gc.Logger = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/data.db", dbDir)), &gc)
	if err != nil {
		return nil
	}
	s.db = db
	// 初始化Id分配器
	s.idAllocator = id.NewIdAllocatorByDB(s.setting.Id, 1000, db)
	// 创建Gin服务
	s.engine = gin.Default()
	s.engine.Use(s.getCors())
	// 创建静态资源
	for _, static := range s.setting.WebConfig.Static {
		s.engine.Static(static[0], static[1])
	}
	for _, static := range s.setting.WebConfig.StaticFile {
		s.engine.StaticFile(static[0], static[1])
	}
	for _, any := range s.setting.WebConfig.Any {
		s.engine.Any(any)
	}

	return s
}

//
// run
//  @Description: 运行服务
//
func (s *Service) run() {
	// 初始化
	err := s.init()
	if err != nil {
		panic(err)
		return
	}

	// 启动服务
	go func() {
		err := s.engine.Run(":" + s.setting.Port)
		if err != nil {
			//f.logAdapter.Fatal("qf run error", err.Error())
			panic(err)
		}
	}()
}

//
// stop
//  @Description: 停止服务
//
func (s *Service) stop() {
	// 执行业务释放
	for _, bll := range s.bllList {
		bll.Stop()
	}
}

//
// RegBll
//  @Description: 注册业务对象
//  @param bll 业务对象
//  @param group 所在组，如果为空则默认为api
//
func (s *Service) RegBll(bll IBll, group string) {
	// 初始化
	pkg, name := s.getBllName(bll)
	bll.setPkg(pkg)
	bll.setName(name)
	bll.setGroup(group)

	// 注册API和路由
	api := ApiMap{}
	bll.RegApi(api)
	router := s.engine.Group(group)
	for kind, routers := range api {
		for relative, handler := range routers {
			path := pkg + "/" + relative
			if kind == EApiKindGetList {
				path = pkg + "s" + "/" + relative
			}
			path = strings.Trim(path, "/")
			router.Handle(kind.HttpMethod(), path, s.context)
			s.apiHandler[fmt.Sprintf("%s:%s/%s", kind.HttpMethod(), group, path)] = handler
		}
	}

	// 注册数据访问层并初始化
	dal := DalMap{}
	bll.RegDal(dal)
	for d, model := range dal {
		// 配置数据库给数据层，并初始化表结构
		d.initDB(s.db, model)
		d.setChild(d)
	}

	// 注册消息
	msg := MessageMap{}
	bll.RegMsg(msg)
	for kind, routers := range msg {
		for relative, handler := range routers {
			sp := strings.Split(relative, ",")
			path := sp[0] + "/" + sp[1]
			if kind == EApiKindGetList {
				path = sp[0] + "s" + "/" + sp[1]
			}
			path = strings.Trim(path, "/")
			s.msgHandler[fmt.Sprintf("%s:%s/%s", kind.HttpMethod(), group, path)] = handler
		}
	}

	//// 注册其他包引用
	//ref := RefMap{}
	//bll.RegRef(ref)
	//for kind, routers := range ref {
	//	for relative, handler := range routers {
	//		sp := strings.Split(relative, ",")
	//		path := sp[0] + "/" + sp[1]
	//		if kind == EApiKindGetList {
	//			path = sp[0] + "s" + "/" + sp[1]
	//		}
	//		path = strings.Trim(path, "/")
	//		s.refHandler[fmt.Sprintf("%s:%s/%s", kind.HttpMethod(), group, path)] = handler
	//	}
	//}

	// 加入到业务列表
	if _, ok := s.bllList[bll.getKey()]; ok == false {
		s.bllList[bll.getKey()] = bll
	} else {
		panic(fmt.Sprintf("%s already exists", bll.getKey()))
	}
}

//
// init
//  @Description: 服务初始化
//  @return error
//
func (s *Service) init() error {
	// 绑定外包引用到业务API
	for _, bll := range s.bllList {
		// 注册其他包引用
		ref := RefMap{}
		ref.bllGroup = bll.getGroup()
		ref.allApis = s.apiHandler
		bll.RegRef(ref)
	}

	// 执行业务初始化
	for _, bll := range s.bllList {
		err := bll.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) getBllName(bll IBll) (string, string) {
	t := reflect.TypeOf(bll).Elem()
	sp := strings.Split(t.PkgPath(), "/")
	return sp[len(sp)-1], t.Name()
}

func (s *Service) context(ctx *gin.Context) {
	url := fmt.Sprintf("%s:%s", ctx.Request.Method, strings.TrimLeft(ctx.FullPath(), "/"))
	if handler, ok := s.apiHandler[url]; ok {
		qfCtx := &Context{
			Time:        time.Now().Local(),
			LoginUser:   LoginUser{}, // TODO
			inputValue:  make([]map[string]interface{}, 0),
			inputSource: "",
			idAllocator: s.idAllocator,
		}
		// 获取body内容
		if body, e := ioutil.ReadAll(ctx.Request.Body); e == nil {
			qfCtx.inputSource = string(body)
			if json.Valid(body) {
				// 如果是json列表
				if strings.HasPrefix(qfCtx.inputSource, "[") &&
					strings.HasSuffix(qfCtx.inputSource, "]") {
					_ = json.Unmarshal(body, &qfCtx.inputValue)
				}
				// 如果是json结构
				if strings.HasPrefix(qfCtx.inputSource, "{") &&
					strings.HasSuffix(qfCtx.inputSource, "}") {
					iv := map[string]interface{}{}
					_ = json.Unmarshal(body, &iv)
					qfCtx.inputValue = append(qfCtx.inputValue, iv)
				}
			} else {
				if qfCtx.inputSource != "" {
					s.returnError(ctx, errors.New("invalid json format"))
					return
				}
			}
		}

		// 获取全部的Query
		for k, v := range ctx.Request.URL.Query() {
			if len(qfCtx.inputValue) == 0 {
				qfCtx.inputValue = append(qfCtx.inputValue, map[string]interface{}{})
			}
			if len(v) > 0 {
				qfCtx.inputValue[0][k] = v[0]
			} else {
				qfCtx.inputValue[0][k] = ""
			}
		}

		// 执行业务方法
		result, err := handler(qfCtx)

		// 返回给前端
		if err != nil {
			s.returnError(ctx, err)
		} else {
			// 判断框架内部业务是否有订阅该消息的
			if mh, ok := s.msgHandler[url]; ok {
				go func() {
					e := mh(qfCtx)
					if e != nil {

					}
				}()
			}

			// TODO：记录日志

			s.returnOk(ctx, result)
		}
	}
}

func (s *Service) returnError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"msg":    err.Error(),
	})
}

func (s *Service) returnOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "success",
		"data":   data,
	})
}

func (s *Service) getCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func BuildContext(ctx *Context, input interface{}) *Context {
	context := &Context{
		Time:        ctx.Time,
		LoginUser:   ctx.LoginUser,
		inputValue:  nil,
		inputSource: "",
		idPer:       ctx.idPer,
		idAllocator: ctx.idAllocator,
	}
	body, _ := json.Marshal(input)
	context.inputSource = string(body)
	if json.Valid(body) {
		// 如果是json列表
		if strings.HasPrefix(context.inputSource, "[") &&
			strings.HasSuffix(context.inputSource, "]") {
			_ = json.Unmarshal(body, &context.inputValue)
		}
		// 如果是json结构
		if strings.HasPrefix(context.inputSource, "{") &&
			strings.HasSuffix(context.inputSource, "}") {
			iv := map[string]interface{}{}
			_ = json.Unmarshal(body, &iv)
			context.inputValue = append(context.inputValue, iv)
		}
	}
	return context
}

//
// buildTableName
//  @Description: 根据结构体，生成对应的数据库表名
//  @param model 结构体
//  @return string 然后表名，规则：包名_结构体名，如果包名和结构体名一致时，则只返回结构体名
//
func buildTableName(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	pName := strings.ToLower(filepath.Base(t.PkgPath()))
	bName := strings.ToLower(t.Name())
	tName := fmt.Sprintf("%s_%s", pName, bName)
	if pName == bName {
		tName = pName
	}
	return tName
}
