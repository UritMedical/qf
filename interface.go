package qf

import (
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
	RegApi(regApi ApiMap)     // 注册需要暴露的API方法
	RegDal(regDal DalMap)     // 注册需要框架初始化的数据访问层对象
	RegMsg(regMsg MessageMap) // 注册需要接收处理的消息
	RegRef(regRef RefMap)     // 注册引用
	Init() error              // 业务自己的初始化方法
	Stop()                    // 业务自己的资源释放方法
	// 框架内部实现的方法
	setPkg(pkg string)                                    // 1
	setName(name string)                                  // 1
	setGroup(group string)                                // 1
	setIConfig(config iConfig)                            // 1
	getKey() string                                       // 1
	getGroup() string                                     // 1
	Debug(content string)                                 // 1
	GetConfig() map[string]interface{}                    // 1
	SetConfig(value map[string]interface{}) (bool, error) // 1
}

//
// BaseBll
//  @Description: 提供业务基础通用方法
//
type BaseBll struct {
	pkg    string
	name   string
	group  string
	config iConfig
}

func (bll *BaseBll) setPkg(pkg string) {
	bll.pkg = pkg
}

func (bll *BaseBll) setName(name string) {
	bll.name = name
}

func (bll *BaseBll) setGroup(group string) {
	bll.group = group
}

func (bll *BaseBll) setIConfig(config iConfig) {
	bll.config = config
}

func (bll *BaseBll) getKey() string {
	return fmt.Sprintf("%s/%s.%s", bll.group, bll.pkg, bll.name)
}

func (bll *BaseBll) getGroup() string {
	return bll.group
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
	return bll.config.GetConfig(bll.getKey())
}

//
// SetConfig
//  @Description:
//  @param value
//  @return bool
//  @return error
//
func (bll *BaseBll) SetConfig(value map[string]interface{}) (bool, error) {
	return bll.config.SetConfig(bll.getKey(), value)
}

//
// IDal
//  @Description: 数据访问层接口
//
type IDal interface {
	DB() *gorm.DB                                                  // 返回数据库对象
	Save(content interface{}) error                                // 执行新增或修改操作
	Delete(id uint64) (bool, error)                                // 执行删除操作
	GetModel(id uint64, dest interface{}) error                    // 根据Id获取单条信息
	GetList(startId uint64, maxCount uint, dest interface{}) error // 根据起始Id和最大数量，获取一组信息
	GetCount(query interface{}, args ...interface{}) int64         // 根据条件获取数量
	CheckExists(id uint64) bool                                    // 检测Id是否存在
	// 框架内部实现的方法
	initDB(db *gorm.DB, model interface{})
	setChild(dal IDal)
}

//
// BaseDal
//  @Description: 基础数据访问方法
//
type BaseDal struct {
	db        *gorm.DB
	dal       IDal
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
func (b *BaseDal) initDB(db *gorm.DB, model interface{}) {
	b.db = db
	// 根据实体名称，生成数据库
	b.tableName = buildTableName(model)
	// 自动生成表
	_ = db.Table(b.tableName).AutoMigrate(model)
}

//
// setChild
//  @Description: 设置子集，用于后期反射
//  @receiver b
//  @param dal
//
func (b *BaseDal) setChild(dal IDal) {
	b.dal = dal
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
func (b *BaseDal) Delete(id uint64) (bool, error) {
	result := b.DB().Delete(&BaseModel{Id: id})
	return result.RowsAffected > 0, result.Error
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

func (msg MessageMap) Reg(pkgName string, kind EApiKind, router string, handler MessageHandler) {
	if msg[kind] == nil {
		msg[kind] = make(map[string]MessageHandler)
	}
	key := fmt.Sprintf("%s,%s", pkgName, router)
	if _, ok := msg[kind][key]; ok == false {
		msg[kind][key] = handler
	} else {
		panic(fmt.Sprintf("【RegMsg】: %s:%s already exists", kind, key))
	}
}

type RefMap struct {
	bllGroup string
	allApis  map[string]ApiHandler
}

func (ref RefMap) Load(pkgName string, kind EApiKind, router string) ApiHandler {
	path := pkgName + "/" + router
	if kind == EApiKindGetList {
		path = pkgName + "s" + "/" + router
	}
	path = strings.Trim(path, "/")
	key := fmt.Sprintf("%s:%s/%s", kind.HttpMethod(), ref.bllGroup, path)
	if api, ok := ref.allApis[key]; ok {
		return api
	}
	return nil
}
