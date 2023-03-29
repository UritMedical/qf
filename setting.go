package qf

import (
	"github.com/UritMedical/qf/util/qconfig"
)

type setting struct {
	Id          uint        `comment:"框架Id，主服务为0"`
	Name        string      `comment:"框架名称，用于网络发现，单体服务可为空"`
	Port        string      `comment:"服务端口"`
	WebConfig   webConfig   `comment:"web配置"`
	UserConfig  userConfig  `comment:"用户配置"`
	GormConfig  gormConfig  `comment:"gorm配置"`
	OtherConfig otherConfig `comment:"其他配置"`
}

type userConfig struct {
	Enabled        byte     `comment:"是否启用用户模块 0否 1是"`
	TokenVerify    string   `comment:"特殊放行的Token密码"`
	TokenWhiteList []string `toml:",multiline" comment:"token白名单，以下路由不进行token验证"`
}

type webConfig struct {
	GinRelease byte       `comment:"是否启动gin的release版本 0否 1是"`
	DefGroup   string     `comment:"路由的默认所在组"`
	Static     [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,root路径"`
	StaticFile [][]string `toml:",multiline" comment:"静态资源配置，格式为：相对路径,文件路径"`
	Any        []string   `toml:",multiline" comment:"特殊路由注册"`
	Mime       [][]string `toml:",multiline" comment:"MIME文件扩展名配置，格式为：文件后缀名,类型"`
	ShortRoute [][]string `toml:",multiline" comment:"短路由配置，格式为：短路由,实际路由（如：/item /dist/#/setting/item）"`
}

type gormConfig struct {
	DBType                 string `comment:"数据库类型：sqlite, sqlserver\n 参数\n sqlite：xxx.db\n sqlserver：ip,db,user,pwd"`
	DBParam                string
	OpenLog                byte   `comment:"是否输出脚本日志 0否 1是"`
	SkipDefaultTransaction byte   `comment:"跳过默认事务 0否 1是"`
	JournalMode            string `comment:"Journal模式\n DELETE：在事务提交后，删除journal文件\n MEMORY：在内存中生成journal文件，不写入磁盘\n WAL：使用WAL（Write-Ahead Logging）模式，将journal记录写入WAL文件中\n OFF：完全关闭journal模式，不记录任何日志消息"`
}

type otherConfig struct {
	JsonDateFormat string `comment:"框架日期的Json格式"`
	JsonTimeFormat string `comment:"框架时间的Json格式"`
}

func (s *setting) Load(path string) {
	// 初始值
	s.Port = "80"
	s.UserConfig = userConfig{
		Enabled:     1,
		TokenVerify: "lis",
		TokenWhiteList: []string{
			"POST:/api/login",
			"POST:/api/user/jwt/reset",
		},
	}
	s.WebConfig = webConfig{
		GinRelease: 0,
		DefGroup:   "api",
		Static: [][]string{
			{"/assets", "./res/assets"},
			{"/js", "./res/js"},
			{"/img", "./res/img"},
			{"/user", "./user"},
		},
		StaticFile: [][]string{
			{"/", "./res/index.html"},
		},
		Any: []string{
			"index.html/*any",
		},
		Mime: [][]string{
			{".js", "text/javascript"},
		},
		ShortRoute: [][]string{},
	}
	s.GormConfig = gormConfig{
		DBType:                 "sqlite",
		DBParam:                "data.db",
		OpenLog:                0,
		SkipDefaultTransaction: 1,
		JournalMode:            "OFF",
	}
	s.OtherConfig = otherConfig{
		JsonDateFormat: "yyyy-MM-dd",
		JsonTimeFormat: "HH:mm:ss",
	}
	_ = qconfig.LoadFromToml(path, s)
}
