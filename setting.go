package qf

import (
	"bytes"
	"github.com/UritMedical/qf/util/io"
	"github.com/pelletier/go-toml/v2"
)

type setting struct {
	Id         uint        `comment:"框架Id，主服务为0"`
	Name       string      `comment:"框架名称，用于网络发现，单体服务可为空"`
	Port       string      `comment:"服务端口"`
	WebConfig  *webConfig  `comment:"web配置"`
	GormConfig *gormConfig `comment:"gorm配置"`
}

type webConfig struct {
	Static     [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,root路径"`
	StaticFile [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,文件路径"`
	Any        []string   `toml:",multiline" comment:"特殊路由注册"`
}

type gormConfig struct {
	OpenLog byte `comment:"是否输出脚本日志 0否 1是"`
}

func (s *setting) Load(path string) {
	data, _ := io.ReadAllBytes(path)
	_ = toml.Unmarshal(data, s)

	changed := false
	if s.Port == "" {
		s.Port = "80"
		changed = true
	}
	if s.WebConfig == nil {
		s.WebConfig = &webConfig{}
		changed = true
	}
	if s.WebConfig.Static == nil {
		s.WebConfig.Static = [][]string{
			{"/assets", "./res/assets"},
			{"/js", "./res/js"},
			{"/img", "./res/img"},
			{"/child", "./child"},
			{"/app1", "./child/app1"},
			{"/app2", "./child/app2"},
		}
		changed = true
	}
	if s.WebConfig.StaticFile == nil {
		s.WebConfig.StaticFile = [][]string{
			{"/", "./res/index.html"},
		}
		changed = true
	}
	if s.WebConfig.Any == nil {
		s.WebConfig.Any = []string{"index.html/*any"}
		changed = true
	}
	if s.GormConfig == nil {
		s.GormConfig = &gormConfig{}
		changed = true
	}
	// 保存
	if changed {
		buf := bytes.Buffer{}
		enc := toml.NewEncoder(&buf)
		enc.SetIndentTables(true)
		_ = enc.Encode(s)
		_ = io.WriteAllBytes(path, buf.Bytes(), false)
	}
}
