package qf

import (
	"fmt"
	"gorm.io/gorm"
)

// Model 基础内容实体，相关业务实体需要集成
type Model struct {
	ID   uint   `gorm:"primarykey"` // 唯一号
	Info string // 完整内容信息
}

// BaseBll 基础业务实现
type BaseBll struct {
	log     ILogAdapter
	setting ISettingAdapter
}

func (b BaseBll) Log() *ILogAdapter {
	return &b.log
}

func (b BaseBll) GetConfig() map[string]interface{} {
	return b.setting.GetConfig()
}

func (b BaseBll) SetConfig(config map[string]interface{}) {
	b.SetConfig(config)
}

// ApiHandler 业务实现
type ApiHandler func(content interface{}) (interface{}, error)

type Api map[string]ApiHandler

func (api Api) Set(kind EKind, router string, handler ApiHandler) {
	key := fmt.Sprintf("%s^%s", kind, router)
	if _, ok := api[key]; ok == false {
		api[key] = handler
	} else {
		panic(fmt.Sprintf("%s already exists", key))
	}
}

type Dal map[IDal]interface{}

func (d *Dal) Set(dal IDal, model interface{}) {

}

type EKind string

var (
	EKindSave     EKind = "Save"
	EKindDelete   EKind = "Delete"
	EKindGetModel EKind = "GetModel"
	EKindGetList  EKind = "GetList"
)

type BaseDal struct {
	db *gorm.DB
}

func (b *BaseDal) setDB(db *gorm.DB) {
	b.db = db
}

func (b *BaseDal) DB() *gorm.DB {
	return b.db
}

func (b *BaseDal) Save(content interface{}) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDal) Delete(content interface{}) (interface{}, error) {
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

type Reference []IBll

func (ref *Reference) Set(bll IBll) {

}
