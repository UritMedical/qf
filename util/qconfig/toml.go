package qconfig

import (
	"bytes"
	"errors"
	"github.com/UritMedical/qf/util/io"
	"github.com/pelletier/go-toml/v2"
	"reflect"
)

//
// LoadFromToml
//  @Description: 通过toml文件加载配置，支持描述输出，需要在指定的属性后面增加如下tag
//                  字段名1 `comment:"描述"`
//                  字段名2 `comment:"描述"`
//  @param path toml
//  @param setting
//  @return error
//
func LoadFromToml(path string, setting interface{}) error {
	value := reflect.ValueOf(setting)
	if value.Kind() != reflect.Ptr {
		return errors.New("the obj 's kind must be ptr")
	}
	// 从文件读取内容
	data, _ := io.ReadAllBytes(path)
	// 反序列化
	err := toml.Unmarshal(data, setting)
	if err != nil {
		return err
	}
	// 然后再序列化后与原始比较，不一致则保存
	after, _ := marshal(setting)
	if len(data) != len(after) {
		return SaveFromToml(path, setting)
	}
	return nil
}

//
// SaveFromToml
//  @Description: 保存到Toml文件
//  @param path
//  @param setting
//  @return error
//
func SaveFromToml(path string, setting interface{}) error {
	buf, err := marshal(setting)
	if err != nil {
		return err
	}
	return io.WriteAllBytes(path, buf, false)
}

func marshal(setting interface{}) ([]byte, error) {
	// 添加tag
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	enc.SetArraysMultiline(true)
	err := enc.Encode(setting)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
