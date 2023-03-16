package qf

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/UritMedical/qf/util/qid"
	"github.com/UritMedical/qf/util/qreflect"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

//
// Context
//  @Description: Api上下文参数
//
type Context struct {
	time      time.Time // 操作时间
	loginUser LoginUser // 登陆用户信息
	// input的原始内容字典
	inputValue  []map[string]interface{}
	inputSource string
	inputFiles  map[string][]*multipart.FileHeader
	// id分配器
	idPer       uint
	idAllocator qid.IIdAllocator
}

//
// NewContext
//  @Description: 生成一个新的上下文
//  @receiver ctx
//  @param input
//  @return *Context
//
func (ctx *Context) NewContext(input interface{}) *Context {
	context := &Context{
		time:        ctx.time,
		loginUser:   ctx.loginUser,
		inputValue:  nil,
		inputSource: "",
		idPer:       ctx.idPer,
		idAllocator: ctx.idAllocator,
	}
	body, _ := json.Marshal(input)
	context.inputSource = string(body)
	if json.Valid(body) {
		// 如果是json列表
		if strings.HasPrefix(context.inputSource, "[") &&
			strings.HasSuffix(context.inputSource, "]") {
			_ = json.Unmarshal(body, &context.inputValue)
		}
		// 如果是json结构
		if strings.HasPrefix(context.inputSource, "{") &&
			strings.HasSuffix(context.inputSource, "}") {
			iv := map[string]interface{}{}
			_ = json.Unmarshal(body, &iv)
			context.inputValue = append(context.inputValue, iv)
		}
	}
	return context
}

//
// NewId
//  @Description: 自动分配一个新ID
//  @param object 表对应的实体对象
//  @return uint64
//
func (ctx *Context) NewId(object interface{}) uint64 {
	return ctx.idAllocator.Next(buildTableName(object))
}

//
// LoginUser
//  @Description: 获取登陆用户信息
//  @return LoginUser
//
func (ctx *Context) LoginUser() LoginUser {
	user := LoginUser{
		UserId:      ctx.loginUser.UserId,
		UserName:    ctx.loginUser.UserName,
		LoginId:     ctx.loginUser.LoginId,
		Departments: map[uint64]struct{ Name string }{},
		token:       ctx.loginUser.token,
		roles:       map[uint64]struct{ Name string }{},
	}
	for id, info := range ctx.loginUser.Departments {
		user.Departments[id] = struct{ Name string }{
			Name: info.Name,
		}
	}
	for id, info := range ctx.loginUser.roles {
		user.roles[id] = struct{ Name string }{
			Name: info.Name,
		}
	}
	return user
}

//
// IsNull
//  @Description: 判断提交的内容是否为空
//  @return bool
//
func (ctx *Context) IsNull() bool {
	if ctx.inputValue == nil || len(ctx.inputValue) == 0 || ctx.inputSource == "" {
		return true
	}
	return false
}

//
// Bind
//  @Description: 绑定到新实体对象
//  @param objectPtr 实体对象指针（必须为指针）
//  @param attachValues 需要附加的值（可以是结构体、字典）
//  @return error
//
func (ctx *Context) Bind(objectPtr interface{}, attachValues ...interface{}) error {
	if objectPtr == nil {
		return errors.New("the object cannot be empty")
	}
	// 创建反射
	ref := qreflect.New(objectPtr)
	// 必须为指针
	if ref.IsPtr() == false {
		return errors.New("the object must be pointer")
	}
	// 追加附加内容到字典
	for _, value := range attachValues {
		r := qreflect.New(value)
		for k, v := range r.ToMap() {
			for i := 0; i < len(ctx.inputValue); i++ {
				ctx.inputValue[i][k] = v
			}
		}
	}
	// 然后根据类型，将字典写入到对象或列表中
	cnt := make([]BaseModel, 0)
	for i := 0; i < len(ctx.inputValue); i++ {
		c := ctx.build(ctx.inputValue[i], ref.ToMap())
		cnt = append(cnt, c)
	}
	// 重新赋值
	return ref.Set(ctx.inputValue, cnt)
}

//
// GetFile
//  @Description: 获取前端上传的文件列表
//  @param key 属性名
//  @return []*multipart.FileHeader
//
func (ctx *Context) GetFile(key string) []*multipart.FileHeader {
	if ctx.inputFiles == nil {
		return nil
	}
	return ctx.inputFiles[key]
}

func (ctx *Context) build(source map[string]interface{}, exclude map[string]interface{}) BaseModel {
	nid := uint64(0)
	if id, ok := source["Id"]; ok {
		v, e := strconv.Atoi(fmt.Sprintf("%v", id))
		if e == nil {
			nid = uint64(v)
		}
	}
	// 从完整的原始input中，去掉实体对象中已经存在的
	finals := map[string]interface{}{}
	for k, v := range source {
		if _, ok := exclude[k]; ok == false {
			finals[k] = v
		}
	}
	info := ""
	if len(finals) > 0 {
		cj, _ := json.Marshal(finals)
		info = string(cj)
	}
	return BaseModel{
		Id:       nid,
		LastTime: ctx.time,
		FullInfo: info,
	}
}

func (ctx *Context) GetJsonValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	nj, _ := json.Marshal(ctx.inputValue[0][propName])
	return string(nj)
}

func (ctx *Context) GetStringValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", ctx.inputValue[0][propName])
}

func (ctx *Context) GetUIntValue(propName string) uint64 {
	num, _ := strconv.Atoi(ctx.GetStringValue(propName))
	return uint64(num)
}

func (ctx *Context) GetId() uint64 {
	return ctx.GetUIntValue("Id")
}
