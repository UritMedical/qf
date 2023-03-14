package qf

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
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
	set(sub IBll, qfGroup, subGroup string, config iConfig) // 将主服务的部分对象设置被基础业务
	key() string                                            // 获取业务唯一编号
	regApi(bind func(key string, handler ApiHandler))       // 框架注册方法
	regMsg(bind func(key string, handler MessageHandler))   // 框架注册方法
	regDal(db *gorm.DB)                                     // 框架注册方法
	regRef(getApi func(key string) ApiHandler)              // 框架注册方法
	Debug(content string)                                   // 调试日志
	GetConfig() map[string]interface{}                      // 获取配置
	SetConfig(value map[string]interface{}) (bool, error)   // 写入配置
}

//
// BaseBll
//  @Description: 提供业务基础通用方法
//
type BaseBll struct {
	pkg      string  // 业务所在的包名
	name     string  // 业务名称
	qfGroup  string  // 框架路径组
	subGroup string  // 自定义路径组
	config   iConfig // 配置接口
	sub      IBll    // 子接口
}

func (bll *BaseBll) set(sub IBll, qfGroup, subGroup string, config iConfig) {
	// 反射子类
	t := reflect.TypeOf(sub).Elem()
	// 初始化
	bll.sub = sub
	bll.pkg = strings.ToLower(t.PkgPath())
	bll.name = strings.ToLower(t.Name())
	bll.qfGroup = strings.ToLower(qfGroup)
	bll.subGroup = strings.ToLower(subGroup)
	bll.config = config
}

func (bll *BaseBll) key() string {
	return fmt.Sprintf("%s.%s", bll.pkg, bll.name)
}

func (bll *BaseBll) regApi(bind func(key string, handler ApiHandler)) {
	api := ApiMap{}
	bll.sub.RegApi(api)
	for kind, routers := range api {
		for relative, handler := range routers {
			bind(bll.buildPathKey(kind, relative), handler)
		}
	}
}

func (bll *BaseBll) regMsg(bind func(key string, handler MessageHandler)) {
	msg := MessageMap{}
	bll.sub.RegMsg(msg)
	for kind, routers := range msg {
		for relative, handler := range routers {
			bind(bll.buildPathKey(kind, relative), handler)
		}
	}
}

func (bll *BaseBll) regDal(db *gorm.DB) {
	dal := DalMap{}
	bll.sub.RegDal(dal)
	for d, model := range dal {
		d.init(db, model)
	}
}

func (bll *BaseBll) regRef(getApi func(key string) ApiHandler) {
	ref := RefMap{
		getApi: func(kind EApiKind, relative string) ApiHandler {
			return getApi(bll.buildPathKey(kind, relative))
		},
	}
	bll.sub.RegRef(ref)
}

func (bll *BaseBll) buildPathKey(kind EApiKind, relative string) string {
	path := fmt.Sprintf("%s/%s/%s", bll.qfGroup, bll.subGroup, relative)
	path = strings.Replace(path, "//", "/", -1)
	path = strings.TrimRight(path, "/")
	return fmt.Sprintf("%s:/%s", kind.HttpMethod(), path)
}

//
// Debug
//  @Description: 日志调试
//  @param content 需要发送的调试信息
//
func (bll *BaseBll) Debug(info string) {
	fmt.Println(info)
}

//
// GetConfig
//  @Description:
//  @return map[string]interface{}
//
func (bll *BaseBll) GetConfig() map[string]interface{} {
	return bll.config.GetConfig(fmt.Sprintf("%s.%s", bll.pkg, bll.name))
}

//
// SetConfig
//  @Description:
//  @param value
//  @return bool
//  @return error
//
func (bll *BaseBll) SetConfig(value map[string]interface{}) (bool, error) {
	return bll.config.SetConfig(fmt.Sprintf("%s.%s", bll.pkg, bll.name), value)
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

//
// BaseDal
//  @Description: 基础数据访问方法
//
type BaseDal struct {
	db        *gorm.DB
	tableName string
}

//
// initDB
//  @Description: 初始化数据库
//  @receiver b
//  @param db
//  @param pkgName
//  @param model
//
func (b *BaseDal) init(db *gorm.DB, model interface{}) {
	b.db = db
	// 根据实体名称，生成数据库
	b.tableName = buildTableName(model)
	// 自动生成表
	err := db.Table(b.tableName).AutoMigrate(model)
	if err != nil {
		panic(fmt.Sprintf("【Gorm】 AutoMigrate %s failed: %s", b.tableName, err.Error()))
	}
}

//
// DB
//  @Description: 返回对应表的数据控制器
//  @return *gorm.DB
//
func (b *BaseDal) DB() *gorm.DB {
	return b.db.Table(b.tableName)
}

//
// Save
//  @Description: 保存内容
//  @param content 包含了内容结构的实体对象
//  @return error 异常
//
func (b *BaseDal) Save(content interface{}) error {
	// 提交
	result := b.DB().Save(content)
	if result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}

//
// Delete
//  @Description: 删除内容
//  @param id 唯一号
//  @return error 异常
//
func (b *BaseDal) Delete(id uint64) error {
	result := b.DB().Delete(&BaseModel{Id: id})
	if result.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("delete failed, id=%d does not exist", id))
	}
	return result.Error
}

//
// GetModel
//  @Description: 获取单条数据
//  @param id 唯一号
//  @param dest 目标实体结构
//  @return error 返回异常
//
func (b *BaseDal) GetModel(id uint64, dest interface{}) error {
	result := b.DB().Where("id = ?", id).Find(dest)
	// 如果异常或者未查询到任何数据
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//
// GetList
//  @Description: 按唯一号区间，获取一组列表
//  @param startId 起始编号
//  @param maxCount 最大获取数
//  @param dest 目标列表
//  @return error 返回异常
//
func (b *BaseDal) GetList(startId uint64, maxCount uint, dest interface{}) error {
	result := b.DB().Limit(int(maxCount)).Offset(int(startId)).Find(dest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//
// GetConditions
//  @Description: 通过自定义条件获取数据
//  @param dest 结构体/列表
//  @param query 条件
//  @param args 条件参数
//  @return error
//
func (b *BaseDal) GetConditions(dest interface{}, query interface{}, args ...interface{}) error {
	result := b.DB().Where(query, args...).Find(dest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//
// GetCount
//  @Description: GetCount
//  @param query 查询条件，如：a = ? and b = ?
//  @param args 条件对应的值
//  @return int64 查询到的记录数
//
func (b *BaseDal) GetCount(query interface{}, args ...interface{}) int64 {
	count := int64(0)
	b.DB().Where(query, args).Count(&count)
	return count
}

//
// CheckExists
//  @Description: 检查内容是否存在
//  @param id 唯一号
//  @return bool true存在 false不存在
//
func (b *BaseDal) CheckExists(id uint64) bool {
	count := int64(0)
	result := b.DB().Where("id = ?", id).Count(&count)
	if count > 0 && result.Error == nil {
		return true
	}
	return false
}

//
// iIdAllocator
//  @Description: Id分配器接口
//
type iIdAllocator interface {
	Next(name string) uint64
}

//
// iConfig
//  @Description: 业务配置文件接口
//
type iConfig interface {
	GetConfig(name string) map[string]interface{}
	SetConfig(name string, value map[string]interface{}) (bool, error)
}

// ApiHandler API指针
type ApiHandler func(ctx *Context) (interface{}, error)

// MessageHandler  消息执行函数指针
type MessageHandler func(ctx *Context) error

// ApiMap API字典
type ApiMap map[EApiKind]map[string]ApiHandler

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

type DalMap map[IDal]interface{}

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

type MessageMap map[EApiKind]map[string]MessageHandler

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

type RefMap struct {
	getApi func(kind EApiKind, router string) ApiHandler
}

func (ref RefMap) Load(kind EApiKind, router string) ApiHandler {
	return ref.getApi(kind, router)
}
