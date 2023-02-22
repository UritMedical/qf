package qf

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	//folderPath string
	db     *gorm.DB
	engine *gin.Engine
	//apis       *Apis
	//
	//actionAdapter  IActionAdapter
	//messageAdapter IMessageAdapter
	//settingAdapter ISettingAdapter
	//logAdapter     ILogAdapter
}

func NewService() *Service {
	s := &Service{
		//apis: &Apis{},
	}
	//// 默认文件夹路径
	//s.folderPath = "."
	//// 创建数据库
	//dbDir := io.CreateDirectory(fmt.Sprintf("%s/db", s.folderPath))
	//db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/data.db", dbDir)), &gorm.Config{})
	//if err != nil {
	//	return nil
	//}
	//s.db = db
	//// 适配器初始化
	//s.actionAdapter = action.NewActionByWebApi()
	//s.settingAdapter = setting.NewSettingByDB(db)
	//s.logAdapter = log.NewLogByFmt()
	//// 创建Gin服务
	//s.engine = gin.Default()
	//s.engine.Use(f.transmit())
	return s
}

func (f *Service) Run() {
	//go func() {
	//	err := f.engine.Run(":80")
	//	if err != nil {
	//		//f.logAdapter.Fatal("qf run error", err.Error())
	//		panic(err)
	//	}
	//}()
}

func (f *Service) Stop() {

}

func (f *Service) RegBll(bll IBll, routerGroup string) {
	//// 基础方法赋值
	//bll.setDB(f.db)
	//bll.setLog(f.logAdapter)
	//bll.setMessage(f.messageAdapter)
	//bll.setSetting(f.settingAdapter)
	//// 创建API方法
	//bll.RegApis(f.apis)
	//// 内容访问器初始化
	//for _, api := range *f.apis {
	//	bll.setContent(content.NewContentByDB(api.Id, f.db))
	//	break
	//}
	//// 执行业务初始化
	//err := bll.Init()
	//if err == nil {
	//	// 创建路由组
	//	router := f.engine.Group(routerGroup)
	//	// 注册路由
	//	for _, api := range *f.apis {
	//		method := ""
	//		switch api.Kind {
	//		case EApiKindSubmit:
	//			method = "POST"
	//		case EApiKindDelete:
	//			method = "DELETE"
	//		case EApiKindGet:
	//			method = "GET"
	//		}
	//		relative := strings.Trim(fmt.Sprintf("%s/%s", api.Id, api.Route), "/")
	//		router.Handle(method, relative)
	//	}
	//	// 注册消息
	//
	//	// 注册引用
	//}
}

func (f *Service) transmit() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		//// 遍历注册的路由，查找对应的路由名字
		//// 然后执行对应的业务方法
		//for k, handler := range f.routerFunc {
		//	if fmt.Sprintf("%s:%s", c.Request.Method, c.FullPath()) == k {
		//		input := content.Content{}
		//		var err error
		//		if c.Request.Method == "GET" {
		//			query := map[string]interface{}{}
		//			err = c.BindQuery(query)
		//			if err == nil {
		//				j, _ := json.Marshal(query)
		//				input.Info = string(j)
		//			}
		//		} else {
		//			err = c.BindJSON(&input)
		//		}
		//		if err != nil {
		//			c.JSON(http.StatusBadRequest, gin.H{
		//				"status": http.StatusBadRequest,
		//				"msg":    err.Error(),
		//			})
		//			return
		//		}
		//		// 执行业务方法
		//		rs, er := handler(input)
		//		// 返回给前端
		//		if er != nil {
		//			c.JSON(http.StatusBadRequest, gin.H{
		//				"status": http.StatusBadRequest,
		//				"msg":    er.Error(),
		//			})
		//		} else {
		//			c.JSON(http.StatusOK, gin.H{
		//				"status": http.StatusOK,
		//				"msg":    "操作成功",
		//				"data":   rs,
		//			})
		//		}
		//	}
		//}
	}
}

//func (f *Service) callback(client mqtt2.Client, m mqtt2.Message) {
//
//}
