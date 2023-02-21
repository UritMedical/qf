package qf

import (
	"gorm.io/gorm"
)

// IBll 业务层接口
type IBll interface {
	Debug(content string)
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
	regApis(apis Apis)
	regDal(dals Dals)
	regReference(refs References)
	init() error
	stop()
}

// IDal 数据层接口
type IDal interface {
	setDB(db *gorm.DB)
	DB() *gorm.DB
	Save(content interface{}) (interface{}, error)
	Delete(content interface{}) (interface{}, error)
	GetModel(content interface{}) (interface{}, error)
	GetList(content interface{}) (interface{}, error)
	BeforeAction(kind EKind, content interface{}) (bool, error)
	AfterAction(kind EKind, content interface{}) (bool, error)
}
