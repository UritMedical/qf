package qf

import (
	"gorm.io/gorm"
	"qf/helper/content"
	"time"
)

// ISetBll 公用方法资源初始化接口
type ISetBll interface {
	setDB(db *gorm.DB)
	setContent(adapter IContentAdapter)
	setSetting(adapter ISettingAdapter)
	setMessage(adapter IMessageAdapter)
	setLog(adapter ILogAdapter)
	Submit(c content.Content) (interface{}, error)
	Delete(contentId uint) (interface{}, error)
	GetModel(contentId uint) (interface{}, error)
	GetList(startTime time.Time, endTime time.Time) (interface{}, error)
}

// IBll 主业务接口
type IBll interface {
	// ISetBll 公共方法配置接口
	ISetBll
	// RegApis 注册需要暴露的Api方法
	RegApis(apis *Apis)
	// RegMessages 注册需要暴露的消息
	RegMessages(messages *Messages)
	// RegReferences 引用其他模块方法的约素及校验
	RegReferences(references *References)
	// Init 初始化
	Init() (err error)
	// Stop 释放
	Stop()

	BeforeApis(kind EApiKind, content content.Content) (interface{}, error)
	AfterApis(kind EApiKind, latest []content.Content, old content.Content) (interface{}, error)
}

// ILogAdapter 日志接口
type ILogAdapter interface {
	Debug(title, content string)
	Info(title, content string)
	Warn(title, content string)
	Error(title, content string)
	Fatal(title, content string)
}

// IMessageAdapter 消息接口，用于模块间的消息通信，默认为eventbus方式
type IMessageAdapter interface {
	Publish(msgId string, payload interface{}) error // 发布消息
}

// IContentAdapter 内容操作接口
type IContentAdapter interface {
	Insert(cnt content.Content) (content.Content, error)
	Update(cnt content.Content) (content.Content, error)
	Save(cnt content.Content) (content.Content, error)
	Delete(id uint) error
	GetModel(id uint) (content.Content, error)
	GetList(startTime, endTime time.Time) ([]content.Content, error)
}

// ISettingAdapter 配置类操作接口
type ISettingAdapter interface {
	Get(id string) string
	Set(id string, value string) (bool, error)
}

type IActionAdapter interface {
	Call(bllId string, method, relative string, content content.Content) (string, error)
}
