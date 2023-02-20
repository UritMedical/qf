package qf

import (
	"fmt"
	"gorm.io/gorm"
	"qf/helper/content"
	"time"
)

type EApiKind string

const (
	EApiKindSubmit EApiKind = "Submit"
	EApiKindDelete EApiKind = "Delete"
	EApiKindGet    EApiKind = "Get"
)

// Api 行为
type Api struct {
	Id        string     // API定义
	Kind      EApiKind   // API种类
	Route     string     // 路由相对路径
	Explain   string     // 帮助说明
	Function  ApiHandler // 执行函数指针
	DemoInput string     // 默认输入结构
}

type ApiRouter struct {
	Id      string
	Explain string
}

type ApiHandler func(content content.Content) (interface{}, error)

type Apis map[string]Api

//func (api Apis) Reg(router ApiRouter, submit ApiHandler, delete ApiHandler, get ApiHandler, demoInput string) {
//	api.CustomReg(router, EApiKindSubmit, "", "保存一条内容，新增或修改", submit, demoInput)
//	api.CustomReg(router, EApiKindDelete, "", "删除一条内容", delete, demoInput)
//	api.CustomReg(router, EApiKindGet, "", "获取单个完整内容", get, demoInput)
//}

func (api Apis) Reg(router ApiRouter, demoInput string) {

}

func (api Apis) CustomReg(router ApiRouter, kind EApiKind, route string, explain string, function ApiHandler, demoInput string) {
	key := fmt.Sprintf("%s^%s^%s", router.Id, kind, route)
	if _, ok := api[key]; ok == false {
		api[key] = Api{
			Id:        router.Id,
			Kind:      kind,
			Route:     route,
			Explain:   fmt.Sprintf("%s（%s）", router.Explain, explain),
			Function:  function,
			DemoInput: demoInput,
		}
	} else {
		panic(fmt.Sprintf("%s already existed", key))
	}
}

type Message struct {
	MessageId string // 消息编号
	Explain   string // 帮助说明
	DemoInput string // 默认输入结构
}

type MessageHandler func(content string) error

type Messages map[string]Message

func (t Messages) Reg(router ApiRouter, addition string, demoInput string) {
	key := fmt.Sprintf("%s^%s", router.Id, addition)
	if _, ok := t[key]; ok == false {
		t[key] = Message{
			MessageId: key,
			Explain:   router.Explain,
			DemoInput: demoInput,
		}
	} else {
		panic(fmt.Sprintf("%s already existed", key))
	}
}

type Reference struct {
	BllId    string   // 业务模块编号
	Kind     EApiKind // API种类
	Relative string   // 路由相对路径
	Explain  string   // 帮助说明
}

type References map[string]Reference

func (ref References) Reg(router ApiRouter, kind EApiKind, relative string) {
	key := fmt.Sprintf("%s^%s", kind, relative)
	if _, ok := ref[key]; ok == false {
		ref[key] = Reference{
			BllId:    router.Id,
			Kind:     kind,
			Relative: relative,
			Explain:  router.Explain,
		}
	} else {
		panic(fmt.Sprintf("%s already existed", key))
	}
}

func (ref References) Init(bll IBll, explain string) {

}

type BaseBll struct {
	DB      *gorm.DB
	Content IContentAdapter
	Log     ILogAdapter     // 日志帮助类
	setting ISettingAdapter // 配置相关帮助类
	message IMessageAdapter // 消息收发帮助类
	api     IActionAdapter  //
}

func (b *BaseBll) setDB(db *gorm.DB) {
	b.DB = db
}

func (b *BaseBll) setContent(adapter IContentAdapter) {
	b.Content = adapter
}

func (b *BaseBll) setSetting(adapter ISettingAdapter) {
	b.setting = adapter
}

func (b *BaseBll) setMessage(adapter IMessageAdapter) {
	b.message = adapter
}

func (b *BaseBll) setLog(adapter ILogAdapter) {
	b.Log = adapter
}

func (b *BaseBll) Submit(c content.Content) (interface{}, error) {
	// 保存内容
	return b.Content.Save(c)
}

func (b *BaseBll) Delete(contentId uint) (interface{}, error) {
	// 删除内容
	err := b.Content.Delete(contentId)
	if err != nil {
		return nil, err
	}
	return contentId, nil
}

func (b *BaseBll) GetModel(contentId uint) (interface{}, error) {
	return b.Content.GetModel(contentId)
}

func (b *BaseBll) GetList(startTime time.Time, endTime time.Time) (interface{}, error) {
	return b.Content.GetList(startTime, endTime)
}

func (b *BaseBll) SendMessage(router ApiRouter, addition string, payload interface{}) error {
	key := fmt.Sprintf("%s^%s", router.Id, addition)
	return b.message.Publish(key, payload)
}

func (b *BaseBll) Call(router ApiRouter, kind EApiKind, relative string, content content.Content) (string, error) {
	method := ""
	switch kind {
	case EApiKindSubmit:
		method = "POST"
	case EApiKindDelete:
		method = "DELETE"
	case EApiKindGet:
		method = "GET"
	}
	return b.api.Call(router.Id, method, relative, content)
}
