package qf

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
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
	j, _ := json.Marshal(model)
	cnt := Content{
		ID:   0,
		Time: time.Now().Local(),
		Info: string(j),
	}
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
	db *gorm.DB
}

func (b *BaseDal) setDB(db *gorm.DB) {
	b.db = db
}

func (b *BaseDal) DB() *gorm.DB {
	return b.db
}

func (b *BaseDal) Save(content interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDal) Delete(content interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDal) GetModel(content interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDal) GetList(content interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
