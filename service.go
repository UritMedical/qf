package qf

import (
	"fmt"
	"github.com/UritMedical/qf/util/launcher"
	"github.com/UritMedical/qf/util/qerror"
	"github.com/UritMedical/qf/util/qid"
	"github.com/UritMedical/qf/util/qio"
	"github.com/UritMedical/qf/util/token"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

var (
	serv       *Service
	regBllFunc func(s *Service)
	stopFunc   func()
)

func doStart() {
	// 创建服务
	serv = newService()
	// 根据配置是否注册用户模块
	if serv.setting.UserConfig.Enabled == 1 {
		serv.userBll = &userBll{}
		serv.RegBll(serv.userBll, "")
	}
	// 注册外部业务
	regBllFunc(serv)
	// 启动服务
	serv.run()
}

func doStop() {
	// 执行外部释放
	if stopFunc != nil {
		stopFunc()
	}
	// 停止服务
	serv.stop()
}

type Service struct {
	folder         string                    // 框架的文件夹路径
	db             *gorm.DB                  // 数据库
	engine         *gin.Engine               // gin
	bllList        map[string]IBll           // 所有创建的业务层对象
	apiHandler     map[string]ApiHandler     // 所有注册的业务API函数指针
	allApis        map[string][]string       // 所有业务包含的API路由
	msgHandler     map[string]MessageHandler // 所有消息执行的函数指针
	errCodes       map[int]string            // 所有故障码字典
	setting        setting                   // 框架配置
	idAllocator    qid.IIdAllocator          // id分配器接口
	loginUser      map[string]LoginUser      // 登陆用户信息
	tokenWhiteList map[string]byte           // token白名单
	userBll        *userBll                  // 用户业务模块
}

//
// newService
//  @Description: 创建框架服务
//  @return *Service 服务对象指针
//
func newService() *Service {
	s := &Service{
		bllList:    map[string]IBll{},
		allApis:    map[string][]string{},
		apiHandler: map[string]ApiHandler{},
		msgHandler: map[string]MessageHandler{},
		loginUser:  map[string]LoginUser{},
		errCodes:   map[int]string{},
		setting:    setting{},
	}
	// 添加通用故障码
	s.errCodes[ErrorCodeParamInvalid] = "无效的参数"
	s.errCodes[ErrorCodePermissionDenied] = "权限不足，拒绝访问"
	s.errCodes[ErrorCodeRecordNotFound] = "未找到记录"
	s.errCodes[ErrorCodeRecordExist] = "记录已经存在"
	s.errCodes[ErrorCodeSaveFailure] = "保存失败"
	s.errCodes[ErrorCodeDeleteFailure] = "删除失败"
	s.errCodes[ErrorCodeOSError] = "系统故障"
	s.errCodes[ErrorCodeUnknown] = "未知故障"
	// 默认文件夹路径
	s.folder = "."
	// 加载配置
	s.setting.Load(fmt.Sprintf("%s/config/config.toml", s.folder))
	s.tokenWhiteList = map[string]byte{}
	for _, t := range s.setting.UserConfig.TokenWhiteList {
		s.tokenWhiteList[t] = 1
	}
	// 初始化gorm
	gc := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		SkipDefaultTransaction: s.setting.GormConfig.SkipDefaultTransaction == 1,
	}
	if s.setting.GormConfig.OpenLog == 1 {
		gc.Logger = logger.Default.LogMode(logger.Info)
	}
	var db *gorm.DB
	var err error
	switch s.setting.GormConfig.DBType {
	case "sqlite":
		// 创建数据库
		dbDir := qio.CreateDirectory(fmt.Sprintf("%s/db", s.folder))
		db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s/%s", dbDir, s.setting.GormConfig.DBParam)), &gc)
		if err != nil {
			return nil
		}
		if s.setting.GormConfig.JournalMode != "" {
			db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", s.setting.GormConfig.JournalMode))
		}
	case "sqlserver":
		sp := strings.Split(s.setting.GormConfig.DBParam, ",")
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", sp[2], sp[3], sp[0], sp[1])
		db, err = gorm.Open(sqlserver.Open(dsn), &gc)
		if err != nil {
			return nil
		}
	}
	if db == nil {
		return nil
	}
	s.db = db
	// 初始化Id分配器
	s.idAllocator = qid.NewIdAllocatorByDB(s.setting.Id, 1001, db)
	// 创建Gin服务
	if s.setting.WebConfig.GinRelease == 1 {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
	}
	s.engine = gin.Default()
	s.engine.Use(s.getCors())
	s.initApiRouter()
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
	// 初始化前端静态文件
	for _, ext := range s.setting.WebConfig.Mime {
		_ = mime.AddExtensionType(ext[0], ext[1])
	}
	// 短路由
	for _, route := range s.setting.WebConfig.ShortRoute {
		s.engine.GET(route[0], func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s://%s/%s", "http", c.Request.Host, route[1]))
		})
	}
	// 其他初始化
	dateFormat = s.setting.OtherConfig.JsonDateFormat
	dateTimeFormat = fmt.Sprintf("%s %s", s.setting.OtherConfig.JsonDateFormat, s.setting.OtherConfig.JsonTimeFormat)

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
		defer qerror.Recover(func(err string) {
			launcher.Exit()
		})
		err := s.engine.Run(":" + s.setting.Port)
		if err != nil {
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
	group = strings.Trim(group, "/")
	// 初始化业务对象
	bll.set(bll, s.setting.WebConfig.DefGroup, group)
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
			// 检测重复接口
			if _, ok := s.apiHandler[key]; ok {
				panic(fmt.Sprintf("【RegApi】: %s:%s already exists", bll.key(), key))
			}
			if _, ok := s.allApis[bll.key()]; ok == false {
				s.allApis[bll.key()] = []string{}
			}
			s.allApis[bll.key()] = append(s.allApis[bll.key()], key)
			s.apiHandler[key] = handler
			sp := strings.Split(key, ":")
			s.engine.Handle(sp[0], sp[1], s.context)
		})
		// 注册数据访问层
		bll.regDal(s.db)
		// 注册消息
		bll.regMsg(func(key string, handler MessageHandler) {
			// 检测重复接口
			if _, ok := s.msgHandler[key]; ok {
				panic(fmt.Sprintf("【RegMsg】: %s: %s already exists", bll.key(), key))
			}
			s.msgHandler[key] = handler
		})
		// 注册异常
		bll.regError(func(code int, err string) {
			// 检测重复接口
			if _, ok := s.errCodes[code]; ok {
				panic(fmt.Sprintf("【RegFault】: %s: %d,%s already exists", bll.key(), code, err))
			}
			s.errCodes[code] = err
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
				panic(fmt.Sprintf("【RegRef】：%s: %s does not exist", bll.key(), key))
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
	defer qerror.Recover(nil)

	url := fmt.Sprintf("%s:%s", ctx.Request.Method, ctx.FullPath())
	if handler, ok := s.apiHandler[url]; ok {
		// 获取Token值
		tkn := ctx.GetHeader("Token")
		bt := ctx.Query("Bi")
		if bt != "" {
			tkn = bt
		}

		// 验证token和权限，返回登陆用户信息
		// 白名单跳过
		login, err := s.verify(tkn, url)
		if err != nil && s.tokenWhiteList[url] == 0 {
			s.returnInvalid(ctx, err)
			return
		}

		// 生成上下文
		qfCtx := &Context{
			loginUser:   login,
			idAllocator: s.idAllocator,
		}
		qfCtx.time.FromTime(time.Now())

		// 解析body
		contentType := ctx.Request.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			// 处理 JSON 数据
			if body, e := ioutil.ReadAll(ctx.Request.Body); e == nil && len(body) > 0 {
				err := qfCtx.loadInput(body)
				if err != nil {
					s.returnError(ctx, Error(ErrorCodeParamInvalid, err.Error()))
					return
				}
			}
		} else if strings.HasPrefix(contentType, "multipart/form-data") {
			// 处理表单数据
			form, err := ctx.MultipartForm()
			if err != nil {
				s.returnError(ctx, Error(ErrorCodeParamInvalid, "invalid form-data"))
				return
			}
			// 将非文件值加入到字典中
			for key, value := range form.Value {
				if len(value) > 0 {
					qfCtx.setInputValue(key, value[0])
				}
			}
			// 获取文件类的值
			qfCtx.inputFiles = form.File
		}

		// 解析Query
		for k, v := range ctx.Request.URL.Query() {
			if len(v) > 0 {
				qfCtx.setInputValue(k, v[0])
			}
		}

		// 重新生成原始内容
		qfCtx.resetSource()

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
						// TODO 日志
					}
				}()
			}
			// 截取登陆接口，获取登陆信息
			if url == fmt.Sprintf("POST:/%s/qf/login", s.setting.WebConfig.DefGroup) {
				tkn = result.(map[string]interface{})["Token"].(string)
				_, err = s.verify(tkn, "")
				if err != nil {
					s.returnInvalid(ctx, err)
					return
				}
			}
			// TODO：记录日志

			s.returnOk(ctx, result)
		}
	}
}

func (s *Service) returnInvalid(ctx *gin.Context, err IError) {
	msg := map[string]interface{}{}
	msg["code"] = err.Code()
	msg["error"] = s.errCodes[err.Code()]
	ctx.JSON(http.StatusUnauthorized, gin.H{
		"status": http.StatusUnauthorized,
		"msg":    msg,
	})
}

func (s *Service) returnError(ctx *gin.Context, err IError) {
	// 记录日志
	qerror.Write(fmt.Sprintf("%s %d %s %s\n", ctx.Request.URL, err.Code(), s.errCodes[err.Code()], err.Error()))

	msg := map[string]interface{}{}
	msg["code"] = err.Code()
	msg["error"] = s.errCodes[err.Code()]
	ctx.JSON(http.StatusBadRequest, gin.H{
		"status": http.StatusBadRequest,
		"msg":    msg,
	})
}

func (s *Service) returnOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"msg":    "success",
		"data":   data,
	})
}

func (s *Service) verify(tokenStr string, url string) (LoginUser, IError) {
	if s.userBll == nil || tokenStr == s.setting.UserConfig.TokenVerify {
		return LoginUser{}, nil
	}
	// 解析token
	claims, err := token.ParseToken(tokenStr)
	if err != nil {
		return LoginUser{}, Error(ErrorCodeTokenInvalid, err.Error())
	}
	t := time.Unix(claims.ExpiresAt, 0)
	if time.Now().After(t) {
		return LoginUser{}, Error(ErrorCodeTokenExpires, err.Error())
	}
	// 如果缓存存在，则通过缓存
	if _, ok := s.loginUser[tokenStr]; ok {
		return s.loginUser[tokenStr], nil
	}
	// 获取用户信息
	user, err := s.userBll.getUserModel(&Context{loginUser: LoginUser{UserId: claims.Id}})
	if err != nil {
		return LoginUser{}, Error(ErrorCodeTokenInvalid, err.Error())
	}

	value := user.(map[string]interface{})
	info := value["Info"].(map[string]interface{})
	login := LoginUser{
		UserId:   uint64(info["Id"].(float64)),
		UserName: info["Name"].(string),
		LoginId:  info["LoginId"].(string),
		roles:    map[uint64]struct{ Name string }{},
	}
	rlist := make([]uint64, 0)
	for _, role := range value["Roles"].([]map[string]interface{}) {
		rid := uint64(role["Id"].(float64))
		login.roles[rid] = struct{ Name string }{
			Name: role["Name"].(string),
		}
		rlist = append(rlist, rid)
	}

	// 获取组织机构
	dps, err := s.userBll.getOrg(login.UserId)
	if err != nil {
		return LoginUser{}, Error(ErrorCodeTokenInvalid, err.Error())
	}
	login.Departments = dps
	if tokenStr != "" {
		s.loginUser[tokenStr] = login
	}

	// 获取api
	login.apis = s.userBll.getUserAllApis(rlist...)

	// 权限验证
	if login.UserId > 2 {
		// 获取用户权限
		if _, exist := login.apis[url]; exist == false {
			e := fmt.Sprintf("the user does not have %s permission", url)
			return login, Error(ErrorCodePermissionDenied, e)
		}
	}

	return login, nil
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

func (s *Service) initApiRouter() {
	router := s.engine.Group(s.setting.WebConfig.DefGroup)
	router.GET("/qf/allApis", func(ctx *gin.Context) {
		apis := map[string][]string{}
		for k, v := range s.allApis {
			name := filepath.Base(k)
			apis[name] = v
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"msg":    "success",
			"data":   apis,
		})
	})
}

//
// Error
//  @Description: 创建故障内容
//  @param code
//  @param err
//  @return IError
//
func Error(code int, err string) IError {
	return errorInfo{
		code:  code,
		error: err,
	}
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
	if strings.HasPrefix(strings.ToLower(t.PkgPath()), "github.com/uritmedical/qf") {
		per = "Qf"
	}
	return fmt.Sprintf("%s%s", per, t.Name())
}
