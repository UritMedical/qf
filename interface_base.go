package qf

import (
	"fmt"
	"gorm.io/gorm"
)

// BaseBll 基础业务实现
type BaseBll struct {
}

func (b BaseBll) Debug(content string) {

}

func (b BaseBll) GetConfig() map[string]interface{} {
	return nil
}

func (b BaseBll) SetConfig(config map[string]interface{}) {

}

// ApiHandler 业务实现
type ApiHandler func(ctx Context) (interface{}, error)

type ApiMap map[string]ApiHandler

func (api ApiMap) Reg(kind EKind, router string, handler ApiHandler) {
	key := fmt.Sprintf("%s^%s", kind, router)
	if _, ok := api[key]; ok == false {
		api[key] = handler
	} else {
		panic(fmt.Sprintf("%s already exists", key))
	}
}

type DalMap map[IDal]interface{}

func (d *DalMap) Reg(dal IDal, model interface{}) {

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

func (b *BaseDal) Save(content interface{}) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (b *BaseDal) Delete(content interface{}) (bool, error) {
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
