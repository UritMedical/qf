package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/helper/config"
	"github.com/UritMedical/qf/helper/id"
	"github.com/UritMedical/qf/util/io"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"net/http"
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
	config      iConfig                   // 配置文件接口
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
	s.setting.Load(fmt.Sprintf("%s/config/config.toml", s.folder))
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
	s.config = config.NewConfigByDB(db)
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
	// 注册
	s.reg()

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
//  @param group 子组路径名
//
func (s *Service) RegBll(bll IBll, group string) {
	// 初始化业务对象
	bll.set(bll, s.setting.UrlGroup, group, s.config)
	// 加入到业务列表
	if _, ok := s.bllList[bll.key()]; ok == false {
		s.bllList[bll.key()] = bll
	} else {
		panic(fmt.Sprintf("%s already exists", bll.key()))
	}
}

func (s *Service) reg() {
	for _, bll := range s.bllList {
		// 注册API和路由
		bll.regApi(func(key string, handler ApiHandler) {
			sp := strings.Split(key, ":")
			s.engine.Handle(sp[0], sp[1], s.context)
			s.apiHandler[key] = handler
		})
		// 注册数据访问层
		bll.regDal(s.db)
		// 注册消息
		bll.regMsg(func(key string, handler MessageHandler) {
			s.msgHandler[key] = handler
		})
	}
}

//
// init
//  @Description: 服务初始化
//  @return error
//
func (s *Service) init() error {
	// 检测注册的消息是否都操作
	for key := range s.msgHandler {
		if _, ok := s.apiHandler[key]; ok == false {
			panic(fmt.Sprintf("【RegMsg】：%s does not exist", key))
		}
	}
	// 绑定包外引用
	for _, bll := range s.bllList {
		bll.regRef(func(key string) ApiHandler {
			if _, ok := s.apiHandler[key]; ok == false {
				panic(fmt.Sprintf("【RegRef】：%s does not exist", key))
			}
			return s.apiHandler[key]
		})
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

func (s *Service) context(ctx *gin.Context) {
	url := fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.FullPath())
	if handler, ok := s.apiHandler[url]; ok {
		qfCtx := &Context{
			time:        time.Now().Local(),
			loginUser:   LoginUser{}, // TODO
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
		time:        ctx.time,
		loginUser:   ctx.loginUser,
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
	per := ""
	// 如果是框架内部业务，则直接增加Qf前缀
	// 反之直接使用实体名称
	if strings.HasPrefix(t.PkgPath(), "github.com/UritMedical/qf") {
		per = "Qf"
	}
	return fmt.Sprintf("%s%s", per, t.Name())
}
