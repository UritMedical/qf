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
	// RegRef
	//  @Description: 注册引用
	//  @param ref
	//
	RegRef(ref RefMap)

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
	// iBll
	//  @Description: 由框架内部实现方法
	//
	iBll
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
	//  @param id
	//  @return error
	//
	Delete(id uint64) (bool, error)

	//
	// GetModel
	//  @Description:
	//  @param id
	//  @param dest
	//  @return error
	//
	GetModel(id uint64, dest interface{}) error

	//
	// GetList
	//  @Description: 按唯一号区间，获取一组列表
	//  @param startId 起始编号
	//  @param maxCount 最大获取数
	//  @param dest 目标列表
	//  @return error 返回异常
	//
	GetList(startId uint64, maxCount uint, dest interface{}) error

	//
	// GetCount
	//  @Description: GetCount
	//  @param query 查询条件，如：a = ? and b = ?
	//  @param args 条件对应的值
	//  @return int64 查询到的记录数
	//
	GetCount(query interface{}, args ...interface{}) int64

	//
	// CheckExists
	//  @Description:
	//  @param id
	//  @return bool
	//
	CheckExists(id uint64) bool

	//
	// iDal
	//  @Description: 由框架内部实现方法
	//
	iDal
}

//
// iBll
//  @Description: 框架内部使用的业务接口方法
//
type iBll interface {
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

//
// iDal
//  @Description: 框架内部使用的数据层接口方法
//
type iDal interface {
	initDB(db *gorm.DB, pkgName string, model interface{})
	setChild(dal IDal)
}
