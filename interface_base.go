package qf

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

//----------------------------------------------------------------

//
// BaseBll
//  @Description: 提供业务基础通用方法
//
type BaseBll struct {
	pkg   string
	name  string
	group string
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

func (bll *BaseBll) getKey() string {
	return fmt.Sprintf("%s/%s.%s", bll.group, bll.pkg, bll.name)
}

//
// Debug
//  @Description: 日志调试
//  @param content 需要发送的调试信息
//
func (bll *BaseBll) Debug(info string) {

}

//
// GetConfig
//  @Description: 获取配置
//  @return map[string]interface{}
//
func (bll *BaseBll) GetConfig() map[string]interface{} {
	return nil
}

//
// SetConfig
//  @Description: 保存配置
//  @param config
//
func (bll *BaseBll) SetConfig(config map[string]interface{}) {

}

//
// BuildContent
//  @Description: 生成新内容
//  @param model
//  @return Content
//
func (bll *BaseBll) BuildContent(model interface{}) Content {
	// 先转一次json
	tj, _ := json.Marshal(model)
	// 然后在反转到内容对象
	cnt := Content{}
	_ = json.Unmarshal(tj, &cnt)

	// 然后重新生成新的Info
	info := map[string]interface{}{}
	_ = json.Unmarshal([]byte(cnt.Info), &info)
	// 反射对象
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		// 通过原始内容
		if field.Kind() == reflect.Struct && field.Type().Name() == "Content" {
			continue
		}
		tag := value.Type().Field(i).Tag.Get("json")
		if tag != "-" {
			info[value.Type().Field(i).Name] = field.Interface()
		}
	}
	nj, _ := json.Marshal(info)
	cnt.Info = string(nj)

	return cnt
}

//----------------------------------------------------------------

// ApiHandler API指针
type ApiHandler func(ctx *Context) (interface{}, error)

// ApiMap API字典
type ApiMap map[EKind]map[string]ApiHandler

//
// Reg
//  @Description: 注册API
//  @param kind 方式
//  @param router 相对路由
//  @param handler 函数指针
//
func (api ApiMap) Reg(kind EKind, router string, handler ApiHandler) {
	if api[kind] == nil {
		api[kind] = make(map[string]ApiHandler)
	}
	if _, ok := api[kind][router]; ok == false {
		api[kind][router] = handler
	} else {
		panic(fmt.Sprintf("%s:%s already exists", kind, router))
	}
}

//----------------------------------------------------------------

type DalMap map[IDal]interface{}

func (d DalMap) Reg(dal IDal, model interface{}) {
	if _, ok := d[dal]; ok == false {
		d[dal] = model
	} else {
		panic(fmt.Sprintf("dal already exists"))
	}
}

type MessageMap map[string]interface{}

type MessageHandler func(ctx *Context) error

func (m *MessageMap) Reg(bll IBll, router string, function MessageHandler) {

}

type BaseDal struct {
	db        *gorm.DB
	dal       IDal
	tableName string
}

func (b *BaseDal) initDB(db *gorm.DB, pkgName string, model interface{}) {
	b.db = db
	// 根据实体名称，生成数据库
	t := reflect.TypeOf(model)
	pName := strings.ToLower(pkgName)
	bName := strings.ToLower(t.Name())
	b.tableName = fmt.Sprintf("%s_%s", pName, bName)
	if pName == bName {
		b.tableName = pName
	}
	// 自动生成表
	_ = db.Table(b.tableName).AutoMigrate(model)
}

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
	// 先执行Before
	err := b.dal.BeforeAction(EKindSave, content)
	if err != nil {
		return err
	}
	// 提交
	result := b.DB().Save(content)
	if result.RowsAffected > 0 {
		// 在执行After
		return b.dal.AfterAction(EKindSave, content)
	}
	return result.Error
}

//
// Delete
//  @Description: 删除内容
//  @param id 唯一号
//  @return error 异常
//
func (b *BaseDal) Delete(id uint) error {
	result := b.DB().Where("id = ?", id).Updates(Content{Id: id, Delete: 1})
	if result.RowsAffected > 0 {
		return nil
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
func (b *BaseDal) GetModel(id uint, dest interface{}) error {
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
func (b *BaseDal) GetList(startId uint, maxCount uint, dest interface{}) error {
	result := b.DB().Limit(int(maxCount)).Offset(int(startId)).Find(dest)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

//
// CheckExists
//  @Description: 检查内容是否存在
//  @param id 唯一号
//  @return bool true存在 false不存在
//
func (b *BaseDal) CheckExists(id uint) bool {
	count := int64(0)
	result := b.DB().Where("id = ?", id).Count(&count)
	if count > 0 && result.Error == nil {
		return true
	}
	return false
}
