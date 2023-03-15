package qconfig

//
// IConfig
//  @Description: 业务配置文件接口
//
type IConfig interface {
	GetConfig(name string) map[string]interface{}
	SetConfig(name string, value map[string]interface{}) (bool, error)
}
