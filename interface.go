package qf

import (
	"gorm.io/gorm"
)

// IBll 业务层接口
type IBll interface {
	Debug(content string)
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
	RegApi(api ApiMap)
	RegDal(dal DalMap)
	RegMsg(msg MessageMap)
	RefBll() []IBll
	Init() error
	Stop()
}

// IDal 数据层接口
type IDal interface {
	setDB(db *gorm.DB)
	DB() *gorm.DB
	Save(content interface{}) (bool, error)
	Delete(content interface{}) (bool, error)
	GetModel(content interface{}) (interface{}, error)
	GetList(content interface{}) (interface{}, error)
	BeforeAction(kind EKind, content interface{}) error
	AfterAction(kind EKind, content interface{}) error
}
