package qf

import (
	"encoding/json"
	"fmt"
	"github.com/UritMedical/qf/util/qid"
	"github.com/UritMedical/qf/util/qreflect"
	"mime/multipart"
	"reflect"
	"strconv"
)

//
// Context
//  @Description: Api上下文参数
//
type Context struct {
	time      DateTime  // 操作时间
	loginUser LoginUser // 登陆用户信息
	// input的原始内容字典
	inputValue  []map[string]interface{}
	inputSource string
	inputFiles  map[string][]*multipart.FileHeader
	// id分配器
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
		idAllocator: ctx.idAllocator,
	}
	// 将body转入到上下文入参
	body, _ := json.Marshal(input)
	_ = context.loadInput(body)
	// 重新生成原始内容
	context.resetSource()
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
	return ctx.loginUser.CopyTo()
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
func (ctx *Context) Bind(objectPtr interface{}, attachValues ...interface{}) IError {
	if objectPtr == nil {
		return Error(ErrorCodeParamInvalid, "the object cannot be empty")
	}

	// 创建反射
	ref := qreflect.New(objectPtr)
	// 必须为指针
	if ref.IsPtr() == false {
		return Error(ErrorCodeParamInvalid, "the object must be pointer")
	}

	// 先用json反转一次
	_ = json.Unmarshal([]byte(ctx.inputSource), objectPtr)

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
	err := ref.Set(ctx.inputValue, cnt)
	if err != nil {
		return Error(ErrorCodeParamInvalid, err.Error())
	}
	return nil
}

//
// LoadFile
//  @Description: 获取前端上传的文件列表
//  @param key form表单的参数名称
//  @return []*multipart.FileHeader
//
func (ctx *Context) LoadFile(key string) []*multipart.FileHeader {
	if ctx.inputFiles == nil {
		return nil
	}
	return ctx.inputFiles[key]
}

//
// GetJsonValue
//  @Description: 获取指定属性值，并返回json格式
//  @param propName
//  @return string
//
func (ctx *Context) GetJsonValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	nj, _ := json.Marshal(ctx.inputValue[0][propName])
	return string(nj)
}

//
// GetStringValue
//  @Description: 获取指定属性值，并返回字符串格式
//  @param propName
//  @return string
//
func (ctx *Context) GetStringValue(propName string) string {
	if len(ctx.inputValue) == 0 {
		return ""
	}
	return fmt.Sprintf("%v", ctx.inputValue[0][propName])
}

//
// GetUIntValue
//  @Description: 获取指定属性值，并返回整形格式
//  @param propName
//  @return uint64
//
func (ctx *Context) GetUIntValue(propName string) uint64 {
	num, _ := strconv.Atoi(ctx.GetStringValue(propName))
	return uint64(num)
}

//
// GetId
//  @Description: 直接获取Id的值
//  @return uint64
//
func (ctx *Context) GetId() uint64 {
	return ctx.GetUIntValue("Id")
}

//-----------------------------------------------------------------------

func (ctx *Context) loadInput(body []byte) error {
	var obj interface{}
	err := json.Unmarshal(body, &obj)
	if err != nil {
		return err
	}
	maps := make([]map[string]interface{}, 0)
	kind := reflect.TypeOf(obj).Kind()
	if kind == reflect.Slice {
		for _, o := range obj.([]interface{}) {
			maps = append(maps, o.(map[string]interface{}))
		}
	} else if kind == reflect.Map || kind == reflect.Struct {
		maps = append(maps, obj.(map[string]interface{}))
	}
	ctx.inputValue = maps
	return nil
}

func (ctx *Context) setInputValue(key string, value interface{}) {
	if len(ctx.inputValue) == 0 {
		ctx.inputValue = append(ctx.inputValue, map[string]interface{}{})
	}
	for i := 0; i < len(ctx.inputValue); i++ {
		ctx.inputValue[i][key] = value
	}
}

func (ctx *Context) resetSource() {
	ctx.inputSource = ""
	if ctx.inputValue != nil {
		if len(ctx.inputValue) == 1 {
			is, err := json.Marshal(ctx.inputValue[0])
			if err == nil {
				ctx.inputSource = string(is)
			}
		} else {
			is, err := json.Marshal(ctx.inputValue)
			if err == nil {
				ctx.inputSource = string(is)
			}
		}
	}
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
