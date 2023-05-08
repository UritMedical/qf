package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

//
// BaseModel
//  @Description: 基础实体对象
//
type BaseModel struct {
	Id      uint64
	Summary string // 摘要
	Info    string // 其他扩展内容
}

type TestModel struct {
	BaseModel

	Name    string
	HisId   *string
	Address string
}

func TestBind(t *testing.T) {
	str := "{" +
		"\"Id\": 0," +
		"\"HisId\": \"12345\"," +
		"\"Name\": \"张三\"," +
		"\"Sex\": \"男\"," +
		"\"Birth\": \"2023-05-05\"," +
		"\"Phone\": \"12345678901\"," +
		"\"IDCard\": \"123456789012345678\"," +
		"\"Address\": \"广西桂林市七星区XX小区\"," +
		"\"SummaryFields\": \"Name,Sex,IDCard\"," +
		"\"InfoFields\": \"\"" +
		"}"
	mp := map[string]interface{}{}
	_ = json.Unmarshal([]byte(str), &mp)

	model := &TestModel{}
	err := Bind(model, mp)
	fmt.Println(err)
	fmt.Println(model)
}
