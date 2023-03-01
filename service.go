package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"net/http"
	"qf/util/io"
	"reflect"
	"strings"
	"time"
)

type Service struct {
	folder     string                // 框架的文件夹路径
	db         *gorm.DB              // 数据库
	engine     *gin.Engine           // gin
	bllList    map[string]IBll       // 所有创建的业务层对象
	apiHandler map[string]ApiHandler // 所有注册的业务API函数指针
	setting    setting               // 框架配置
}

//
// NewService
//  @Description: 创建框架服务
//  @return *Service 服务对象指针
//
func NewService() *Service {
	s := &Service{
		bllList:    map[string]IBll{},
		apiHandler: map[string]ApiHandler{},
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
// Run
//  @Description: 运行服务
//
func (s *Service) Run() {
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
// Stop
//  @Description: 停止服务
//
func (s *Service) Stop() {
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
			if kind == EKindGetList {
				path = pkg + "s" + "/" + relative
			}
			router.Handle(kind.HttpMethod(), path, s.context)
			s.apiHandler[fmt.Sprintf("%s:%s/%s", kind.HttpMethod(), group, path)] = handler
		}
	}

	// 注册数据访问层并初始化
	dal := DalMap{}
	bll.RegDal(dal)
	for d, model := range dal {
		// 配置数据库给数据层，并初始化表结构
		d.initDB(s.db, pkg, model)
		d.setChild(d)
	}

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
	// 给所有引用的第三方业务赋值
	for _, bll := range s.bllList {
		refs := bll.RefBll()
		for i := 0; i < len(refs); i++ {
			if b, ok := s.bllList[refs[i].getKey()]; ok {
				refs[i] = b
			} else {
				panic("not found")
			}
		}
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
			Time:      time.Now().Local(),
			UserId:    1,
			UserName:  "暂时写死测试用",
			jsonValue: map[string]interface{}{},
		}
		// 获取body内容
		if body, e := ioutil.ReadAll(ctx.Request.Body); e == nil {
			qfCtx.stringValue = string(body)
			if json.Valid(body) {
				// 如果是json格式，则转为字典
				if strings.HasPrefix(qfCtx.stringValue, "{") &&
					strings.HasSuffix(qfCtx.stringValue, "}") {
					_ = json.Unmarshal(body, &qfCtx.jsonValue)
				}
			} else {
				if qfCtx.stringValue != "" {
					s.returnError(ctx, errors.New("invalid json format"))
					return
				}
			}
		}

		// 获取全部的Query
		for k, v := range ctx.Request.URL.Query() {
			if len(v) > 0 {
				qfCtx.jsonValue[k] = v[0]
			} else {
				qfCtx.jsonValue[k] = ""
			}
		}

		// 执行业务方法
		rs, err := handler(qfCtx)

		// TODO：记录日志

		// 返回给前端
		if err != nil {
			s.returnError(ctx, err)
		} else {
			s.returnOk(ctx, rs)
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

//
// BuildContext
//  @Description: 生成上下位对象
//  @param value
//  @return Context
//
func BuildContext(value map[string]interface{}) Context {
	return Context{}
}
