package qf

import (
	"gorm.io/gorm"
)

// IBll 业务层接口
type IBll interface {
	Log() *ILogAdapter
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
	SetApis(api Api)
	SetDal(dal Dal)
	SetReference(ref Reference)
	Init() error
	Stop()
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

// ILogAdapter 日志接口
type ILogAdapter interface {
	Debug(title, content string)
	Info(title, content string)
	Warn(title, content string)
	Error(title, content string)
	Fatal(title, content string)
}

// ISettingAdapter 配置接口
type ISettingAdapter interface {
	GetConfig() map[string]interface{}
	SetConfig(config map[string]interface{})
}
