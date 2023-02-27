package qf

import (
	"gorm.io/gorm"
)

//
// IBll
//  @Description: 通用业务接口方法
//
type IBll interface {
	//
	// RegApi
	//  @Description: 注册需要暴露的API方法
	//  @param api api字典
	//
	RegApi(api ApiMap)

	//
	// RegDal
	//  @Description: 注册需要框架初始化的数据访问层对象
	//  @param dal 数据层字典
	//
	RegDal(dal DalMap)

	//
	// RegMsg
	//  @Description: 注册需要接收处理的消息
	//  @param msg 消息字典
	//
	RegMsg(msg MessageMap)

	//
	// RefBll
	//  @Description: 提交需要引用的第三方业务，由框架进行初始化
	//  @return []IBll
	//
	RefBll() []IBll

	//
	// Init
	//  @Description: 业务自己的初始化方法
	//  @return error
	//
	Init() error

	//
	// Stop
	//  @Description: 业务自己的资源释放方法
	//
	Stop()

	//
	// IQFBll
	//  @Description: 由框架内部实现方法
	//
	IQFBll
}

// IDal 数据层接口
type IDal interface {
	//
	// DB
	//  @Description: 返回数据库对象
	//  @return *gorm.DB
	//
	DB() *gorm.DB

	//
	// Save
	//  @Description: 执行新增或修改操作
	//  @param content
	//  @return error
	//
	Save(content interface{}) error

	//
	// Delete
	//  @Description: 执行删除操作
	//  @param content
	//  @return error
	//
	Delete(content interface{}) error

	//
	// GetModel
	//  @Description: 获取单个内容信息
	//  @param content
	//  @return interface{}
	//  @return error
	//
	GetModel(content interface{}) (interface{}, error)

	//
	// GetList
	//  @Description:
	//  @param startId
	//  @param maxCount
	//  @return interface{}
	//  @return error
	//
	GetList(startId uint, maxCount int) (interface{}, error)

	//
	// BeforeAction
	//  @Description: 内置增删改查执行前触发
	//  @param kind
	//  @param content
	//  @return error
	//
	BeforeAction(kind EKind, content interface{}) error

	//
	// AfterAction
	//  @Description: 内置增删改查执行后触发
	//  @param kind
	//  @param content
	//  @return error
	//
	AfterAction(kind EKind, content interface{}) error

	//
	// IQFDal
	//  @Description: 由框架内部实现方法
	//
	IQFDal
}

//
// IQFBll
//  @Description: 框架内部使用的业务接口方法
//
type IQFBll interface {
	setPkg(pkg string)
	setName(name string)
	setGroup(group string)
	getKey() string

	//
	// Debug
	//  @Description: 调试日志输出
	//  @param content
	//
	Debug(content string)

	//
	// GetConfig
	//  @Description: 获取配置
	//  @return map[string]interface{}
	//
	GetConfig() map[string]interface{}

	//
	// SetConfig
	//  @Description: 写入配置
	//  @param config
	//
	SetConfig(config map[string]interface{})
}

type IQFDal interface {
	//
	// SetDB
	//  @Description: 设置数据库
	//  @param db
	//
	setDB(db *gorm.DB)
}
