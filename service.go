package qf

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"qf/util/io"
	"reflect"
	"strings"
)

type Service struct {
	folder     string                // 框架的文件夹路径
	db         *gorm.DB              // 数据库
	engine     *gin.Engine           // gin
	bllList    map[string]IBll       // 所有创建的业务层对象
	apiHandler map[string]ApiHandler // 所有注册的业务API函数指针
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
	}
	// 默认文件夹路径
	s.folder = "."
	// 创建数据库
	dbDir := io.CreateDirectory(fmt.Sprintf("%s/db", s.folder))
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s/data.db", dbDir)), &gorm.Config{})
	if err != nil {
		return nil
	}
	s.db = db
	// 创建Gin服务
	s.engine = gin.Default()
	//s.engine.Use(f.transmit())
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
	}

	// 启动服务
	go func() {
		err := s.engine.Run(":80")
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
		// 根据实体名称，生成数据库
		t := reflect.TypeOf(model)
		db := s.db.Table(fmt.Sprintf("%s_%s", strings.ToLower(pkg), strings.ToLower(t.Name())))
		_ = db.AutoMigrate(model)
		// 初始化数据层
		d.setDB(db)
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

		bll.Init()
	}
	// 给所有业务的数据访问对象初始化

	return nil
}

func (s *Service) getBllName(bll IBll) (string, string) {
	t := reflect.TypeOf(bll).Elem()
	sp := strings.Split(t.PkgPath(), "/")
	return sp[len(sp)-1], t.Name()
}

func (s *Service) getModelName(model interface{}) string {
	t := reflect.TypeOf(model)
	return t.Name()
}

func (s *Service) transmit() gin.HandlerFunc {
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

func (s *Service) context(ctx *gin.Context) {

}

//func (f *Service) callback(client mqtt2.Client, m mqtt2.Message) {
//
//}
