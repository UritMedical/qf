package qf

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

//
// BaseBll
//  @Description: 提供业务基础通用方法
//
type BaseBll struct {
	pkg      string // 业务所在的包名
	name     string // 业务名称
	qfGroup  string // 框架路径组
	subGroup string // 自定义路径组
	sub      IBll   // 子接口
}

func (bll *BaseBll) set(sub IBll, qfGroup, subGroup string) {
	// 反射子类
	t := reflect.TypeOf(sub).Elem()
	// 初始化
	bll.sub = sub
	bll.pkg = strings.ToLower(t.PkgPath())
	bll.name = strings.ToLower(t.Name())
	bll.qfGroup = strings.ToLower(qfGroup)
	bll.subGroup = strings.ToLower(subGroup)
}

func (bll *BaseBll) key() string {
	return fmt.Sprintf("%s.%s", bll.pkg, bll.name)
}

func (bll *BaseBll) regApi(bind func(key string, handler ApiHandler)) {
	api := ApiMap{
		bllName: bll.key(),
		dict:    map[EApiKind]map[string]ApiHandler{},
	}
	bll.sub.RegApi(api)
	for kind, routers := range api.dict {
		for relative, handler := range routers {
			bind(bll.buildPathKey(kind, relative), handler)
		}
	}
}

func (bll *BaseBll) regMsg(bind func(key string, handler MessageHandler)) {
	msg := MessageMap{
		bllName: bll.key(),
		dict:    map[EApiKind]map[string]MessageHandler{},
	}
	bll.sub.RegMsg(msg)
	for kind, routers := range msg.dict {
		for relative, handler := range routers {
			bind(bll.buildPathKey(kind, relative), handler)
		}
	}
}

func (bll *BaseBll) regDal(db *gorm.DB) {
	dal := DalMap{
		bllName: bll.key(),
		dict:    map[IDal]interface{}{},
	}
	bll.sub.RegDal(dal)
	for d, model := range dal.dict {
		d.init(db, model)
	}
}

func (bll *BaseBll) regError(bind func(code int, error string)) {
	err := FaultMap{
		bllName: bll.key(),
		dict:    map[int]string{},
	}
	bll.sub.RegFault(err)
	for code, msg := range err.dict {
		bind(code, msg)
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

////
//// GetConfig
////  @Description:
////  @return map[string]interface{}
////
//func (bll *BaseBll) GetConfig() map[string]interface{} {
//	return bll.config.GetConfig(fmt.Sprintf("%s.%s", bll.pkg, bll.name))
//}
//
////
//// SetConfig
////  @Description:
////  @param value
////  @return bool
////  @return error
////
//func (bll *BaseBll) SetConfig(value map[string]interface{}) (bool, error) {
//	return bll.config.SetConfig(fmt.Sprintf("%s.%s", bll.pkg, bll.name), value)
//}

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
	if model != nil {
		b.tableName = buildTableName(model)
		// 自动生成表
		err := db.Table(b.tableName).AutoMigrate(model)
		if err != nil {
			panic(fmt.Sprintf("【Gorm】 AutoMigrate %s failed: %s", b.tableName, err.Error()))
		}
	}
}

//
// DB
//  @Description: 返回对应表的数据控制器
//  @return *gorm.DB
//
func (b *BaseDal) DB() *gorm.DB {
	if b.tableName == "" {
		return b.db
	}
	return b.db.Table(b.tableName)
}

//
// Create
//  @Description: 新增内容
//  @param content 包含了内容结构的实体对象
//  @return IError 异常
//
func (b *BaseDal) Create(content interface{}) IError {
	// 提交
	result := b.DB().Create(content)
	if result.RowsAffected > 0 {
		return nil
	}
	if result.Error != nil {
		return Error(ErrorCodeSaveFailure, result.Error.Error())
	}
	return nil
}

//
// Save
//  @Description: 保存内容（存在更新、不存在则新增）
//  @param content 包含了内容结构的实体对象
//  @return IError 异常
//
func (b *BaseDal) Save(content interface{}) IError {
	// 提交
	result := b.DB().Save(content)
	if result.RowsAffected > 0 {
		return nil
	}
	if result.Error != nil {
		return Error(ErrorCodeSaveFailure, result.Error.Error())
	}
	return nil
}

//
// Delete
//  @Description: 删除内容
//  @param id 唯一号
//  @return IError 异常
//
func (b *BaseDal) Delete(id uint64) IError {
	result := b.DB().Delete(&BaseModel{Id: id})
	if result.RowsAffected == 0 {
		return Error(ErrorCodeRecordNotFound, fmt.Sprintf("delete failed, id=%d does not exist", id))
	}
	if result.Error != nil {
		return Error(ErrorCodeDeleteFailure, result.Error.Error())
	}
	return nil
}

//
// GetModel
//  @Description: 获取单条数据
//  @param id 唯一号
//  @param dest 目标实体结构
//  @return IError 返回异常
//
func (b *BaseDal) GetModel(id uint64, dest interface{}) IError {
	result := b.DB().Where("id = ?", id).Find(dest)
	// 如果异常或者未查询到任何数据
	if result.Error != nil {
		return Error(ErrorCodeRecordNotFound, result.Error.Error())
	}
	return nil
}

//
// GetList
//  @Description: 按唯一号区间，获取一组列表
//  @param startId 起始编号
//  @param maxCount 最大获取数
//  @param dest 目标列表
//  @return IError 返回异常
//
func (b *BaseDal) GetList(startId uint64, maxCount uint, dest interface{}) IError {
	result := b.DB().Limit(int(maxCount)).Offset(int(startId)).Find(dest)
	if result.Error != nil {
		return Error(ErrorCodeRecordNotFound, result.Error.Error())
	}
	return nil
}

//
// GetConditions
//  @Description: 通过自定义条件获取数据
//  @param dest 结构体/列表
//  @param query 条件
//  @param args 条件参数
//  @return IError
//
func (b *BaseDal) GetConditions(dest interface{}, query interface{}, args ...interface{}) IError {
	result := b.DB().Where(query, args...).Find(dest)
	if result.Error != nil {
		return Error(ErrorCodeRecordNotFound, result.Error.Error())
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
