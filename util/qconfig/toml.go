package qconfig

import (
	"bytes"
	"errors"
	"github.com/UritMedical/qf/util/io"
	"github.com/pelletier/go-toml/v2"
	"reflect"
	"strings"
)

//
// LoadFromToml
//  @Description: 通过toml文件加载配置，支持描述输出，需要在指定的属性后面增加如下tag
//                  字段名1 `comment:"描述"`
//                  字段名2 `comment:"描述"`
//  @param path 配置文件路径 xxx.toml
//  @param config 配置结构体
//  @return error
//
func LoadFromToml(path string, config interface{}) error {
	value := reflect.ValueOf(config)
	if value.Kind() != reflect.Ptr {
		return errors.New("the obj 's kind must be ptr")
	}
	// 从文件读取内容
	data, _ := io.ReadAllBytes(path)
	// 反序列化
	err := toml.Unmarshal(data, config)
	if err != nil {
		return err
	}
	// 然后再序列化后与原始比较，不一致则保存
	after, _ := marshal(config)
	if len(data) != len(after) {
		return SaveFromToml(path, config)
	}
	return nil
}

//
// SaveFromToml
//  @Description: 保存到Toml文件
//  @param path 配置文件路径 xxx.toml
//  @param config 配置结构体
//  @return error
//
func SaveFromToml(path string, config interface{}) error {
	buf, err := marshal(config)
	if err != nil {
		return err
	}
	return io.WriteAllBytes(path, buf, false)
}

func marshal(setting interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := toml.NewEncoder(&buf)
	enc.SetIndentTables(true)
	if existMultilineTag(reflect.TypeOf(setting)) == false {
		enc.SetArraysMultiline(true)
	}
	err := enc.Encode(setting)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func existMultilineTag(tp reflect.Type) bool {
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if field.Type.Kind() == reflect.Struct {
			ok := existMultilineTag(field.Type)
			if ok {
				return true
			}
		}
		if tag, ok := tp.Field(i).Tag.Lookup("toml"); ok {
			if strings.Contains(tag, "multiline") {
				return true
			}
		}
	}
	return false
}
