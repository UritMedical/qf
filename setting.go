package qf

import (
	"bytes"
	"github.com/UritMedical/qf/util/io"
	"github.com/pelletier/go-toml/v2"
	"strings"
)

type setting struct {
	Id          uint         `comment:"框架Id，主服务为0"`
	Name        string       `comment:"框架名称，用于网络发现，单体服务可为空"`
	Port        string       `comment:"服务端口"`
	UserConfig  *userConfig  `comment:"用户配置"`
	WebConfig   *webConfig   `comment:"web配置"`
	GormConfig  *gormConfig  `comment:"gorm配置"`
	OtherConfig *otherConfig `comment:"其他配置"`
}

type userConfig struct {
	Enabled        byte     `comment:"是否启用用户模块 0否 1是"`
	TokenVerify    byte     `comment:"是否启用全局token验证 0否 1是"`
	TokenWhiteList []string `toml:",multiline" comment:"token白名单路由列表，以下路由不进行token验证"`
}

type webConfig struct {
	DefGroup   string     `comment:"路由的默认所在组"`
	Static     [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,root路径"`
	StaticFile [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,文件路径"`
	Any        []string   `toml:",multiline" comment:"特殊路由注册"`
}

type gormConfig struct {
	DBName                 string `comment:"默认数据库名称"`
	OpenLog                byte   `comment:"是否输出脚本日志 0否 1是"`
	SkipDefaultTransaction byte   `comment:"跳过默认事务 0否 1是"`
	JournalMode            string `comment:"跳过默认事务\n DELETE：在事务提交后，删除journal文件\n MEMORY：在内存中生成journal文件，不写入磁盘\n WAL：使用WAL（Write-Ahead Logging）模式，将journal记录写入WAL文件中\n OFF：完全关闭journal模式，不记录任何日志消息"`
}

type otherConfig struct {
	JsonDateFormat string `comment:"框架日期的Json格式"`
	JsonTimeFormat string `comment:"框架时间的Json格式"`
}

func (s *setting) Load(path string) {
	data, _ := io.ReadAllBytes(path)
	content := string(data)
	_ = toml.Unmarshal(data, s)

	changed := false
	if s.Port == "" {
		s.Port = "80"
		changed = true
	}
	if s.UserConfig == nil {
		s.UserConfig = &userConfig{
			Enabled:     1,
			TokenVerify: 1,
			TokenWhiteList: []string{
				"POST:/api/login",
				"POST:/api/user/parseToken",
				"POST:/api/user/jwt/reset",
			},
		}
		changed = true
	}
	if s.WebConfig == nil {
		s.WebConfig = &webConfig{}
	}
	if strings.Contains(content, "DefGroup") == false {
		s.WebConfig.DefGroup = "api"
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
	if s.GormConfig.DBName == "" {
		s.GormConfig.DBName = "data"
		changed = true
	}
	if strings.Contains(content, "SkipDefaultTransaction") == false {
		s.GormConfig.SkipDefaultTransaction = 1
		changed = true
	}
	if strings.Contains(content, "JournalMode") == false {
		s.GormConfig.JournalMode = "OFF"
		changed = true
	}
	if s.OtherConfig == nil {
		s.OtherConfig = &otherConfig{}
		changed = true
	}
	if strings.Contains(content, "JsonDateFormat") == false {
		s.OtherConfig.JsonDateFormat = "yyyy-MM-dd"
		changed = true
	}
	if strings.Contains(content, "JsonTimeFormat") == false {
		s.OtherConfig.JsonTimeFormat = "HH:mm:ss"
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
