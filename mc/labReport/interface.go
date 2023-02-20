/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/20 9:10
 */

package labReport

import (
	"gorm.io/gorm"
	"qf/helper/content"
	"time"
)

type IBll interface {
	GetConfig() map[string]interface{}
	SetConfig(map[string]interface{})

	SendMessage(topic string, msg map[string]interface{}) error
	SendLog(level string, log string)
	OnMessage(topic string, message interface{})
	GetDals() map[string]IDal

	//SetAPIs 暴露api
	SetAPIs(apis map[string]ApiHandler) error

	//GetAPIs 获取外部apis引用
	GetAPIs() map[string]ApiHandler
	GetMessages() map[string]bool
}
type IDal interface {
	//Parent() 范例里推荐实现 获取其bll对象  以解决SendLog SendMessage和与其他dal之间互相访问的问题

	GetDb() *gorm.DB                                                                   //base
	Save(cnt content.Content) (content.Content, error)                                 //base
	Delete(id uint) error                                                              //base
	GetModel(id uint) (content.Content, error)                                         //base
	GetList(startTime, endTime time.Time, haveDeleted bool) ([]content.Content, error) //base
	BeforeSave(content *content.Content) error
	AfterSave(content *content.Content) error
}
type ApiHandler func(map[string]interface{}) (interface{}, error)
type MessageHandler func(topic string, msg map[string]interface{}) error
