package qf

import (
	"fmt"
	"github.com/UritMedical/qf/util/qconfig"
	"gorm.io/gorm"
	"reflect"
	"time"
)

//
// IBll
//  @Description: 通用业务接口方法
//
type IBll interface {
	RegApi(a ApiMap)     // 注册需要暴露的API方法
	RegDal(d DalMap)     // 注册需要框架初始化的数据访问层对象
	RegMsg(m MessageMap) // 注册需要接收处理的消息
	RegRef(r RefMap)     // 注册引用
	Init() error         // 业务自己的初始化方法
	Stop()               // 业务自己的资源释放方法
	// 框架内部实现的方法
	set(sub IBll, qfGroup, subGroup string, config qconfig.IConfig) // 将主服务的部分对象设置被基础业务
	key() string                                                    // 获取业务唯一编号
	regApi(bind func(key string, handler ApiHandler))               // 框架注册方法
	regMsg(bind func(key string, handler MessageHandler))           // 框架注册方法
	regDal(db *gorm.DB)                                             // 框架注册方法
	regRef(getApi func(key string) ApiHandler)                      // 框架注册方法
	Debug(content string)                                           // 调试日志
	GetConfig() map[string]interface{}                              // 获取配置
	SetConfig(value map[string]interface{}) (bool, error)           // 写入配置
}

//
// IDal
//  @Description: 数据访问层接口
//
type IDal interface {
	DB() *gorm.DB                                                  // 返回数据库对象
	Save(content interface{}) error                                // 执行新增或修改操作
	Delete(id uint64) error                                        // 执行删除操作
	GetModel(id uint64, dest interface{}) error                    // 根据Id获取单条信息
	GetList(startId uint64, maxCount uint, dest interface{}) error // 根据起始Id和最大数量，获取一组信息
	GetCount(query interface{}, args ...interface{}) int64         // 根据条件获取数量
	CheckExists(id uint64) bool                                    // 检测Id是否存在
	// 框架内部实现的方法
	init(db *gorm.DB, model interface{})
}

// EApiKind 行为类别
type EApiKind string

var (
	EApiKindSave     EApiKind = "Save"     // 新增或修改
	EApiKindDelete   EApiKind = "Delete"   // 删除
	EApiKindGetModel EApiKind = "GetModel" // 获取单条
	EApiKindGetList  EApiKind = "GetList"  // 获取多条
)

//
// HttpMethod
//  @Description: 返回Http方式名称
//  @return string
//
func (kind EApiKind) HttpMethod() string {
	if kind == EApiKindSave {
		return "POST"
	}
	if kind == EApiKindDelete {
		return "DELETE"
	}
	return "GET"
}

// 注册相关结构定义
type (
	ApiMap     map[EApiKind]map[string]ApiHandler     // API字典
	DalMap     map[IDal]interface{}                   // 数据访问字典
	MessageMap map[EApiKind]map[string]MessageHandler // 消息字典
	RefMap     struct {
		getApi func(kind EApiKind, router string) ApiHandler
	} // 外部引用字典
	ApiHandler     func(ctx *Context) (interface{}, error) // API函数指针
	MessageHandler func(ctx *Context) error                // 消息函数指针
)

//
// Reg
//  @Description: 注册对外暴露的方法
//  @param kind 类型
//  @param router 路由相对路径
//  @param handler 函数指针
//
func (api ApiMap) Reg(kind EApiKind, router string, handler ApiHandler) {
	if api[kind] == nil {
		api[kind] = make(map[string]ApiHandler)
	}
	if _, ok := api[kind][router]; ok == false {
		api[kind][router] = handler
	} else {
		panic(fmt.Sprintf("api.reg: %s:%s already exists", kind, router))
	}
}

//
// Reg
//  @Description: 注册数据访问层，并初始化数据库
//  @param iDal 访问层对象
//  @param model 数据库表实体
//
func (d DalMap) Reg(iDal IDal, model interface{}) {
	t := reflect.TypeOf(iDal)
	if t.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("【RegDal】: %s/%s this model must be of type pointer", t.PkgPath(), t.Name()))
	}
	t = t.Elem()
	v := reflect.ValueOf(iDal)
	if v.IsNil() {
		panic(fmt.Sprintf("【RegDal】: %s/%s has not been initialized", t.PkgPath(), t.Name()))
	}
	if _, ok := d[iDal]; ok == false {
		d[iDal] = model
	} else {

		panic(fmt.Sprintf("【RegDal】: %s/%s already exists", t.PkgPath(), t.Name()))
	}
}

//
// Reg
//  @Description: 注册消息方法
//  @param kind 类型
//  @param router 路由相对路径
//  @param handler 函数指针
//
func (msg MessageMap) Reg(kind EApiKind, router string, handler MessageHandler) {
	if msg[kind] == nil {
		msg[kind] = make(map[string]MessageHandler)
	}
	if _, ok := msg[kind][router]; ok == false {
		msg[kind][router] = handler
	} else {
		panic(fmt.Sprintf("【RegMsg】: %s:%s already exists", kind, router))
	}
}

//
// Load
//  @Description: 加载外部方法
//  @param kind 类型
//  @param router 路由相对路径
//  @return ApiHandler 方法函数指针
//
func (ref RefMap) Load(kind EApiKind, router string) ApiHandler {
	return ref.getApi(kind, router)
}

//
// BaseModel
//  @Description: 基础实体对象
//
type BaseModel struct {
	Id       uint64    `gorm:"primaryKey"` // 唯一号
	LastTime time.Time `gorm:"index"`      // 最后操作时间时间
	FullInfo string    // 内容
}

//
// LoginUser
//  @Description: 登陆用户信息
//
type LoginUser struct {
	UserId      uint64 // 登陆用户唯一号
	UserName    string // 登陆用户名字
	LoginId     string // 登陆用户账号
	Departments map[uint64]struct {
		Name string
	} // 所属部门列表
	token string // 登陆的token信息
	roles map[uint64]struct {
		Name string
	} // 角色列表
}
