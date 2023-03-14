package config

import (
	"encoding/json"
	"gorm.io/gorm"
)

type ByDB struct {
	name string
	db   *gorm.DB
}

type config struct {
	Name  string `gorm:"primaryKey"`
	Value string
}

func NewConfigByDB(db *gorm.DB) *ByDB {
	name := "QfBllConfig"
	_ = db.Table(name).AutoMigrate(config{})
	return &ByDB{
		name: name,
		db:   db,
	}
}

//
// GetConfig
//  @Description: 获取配置
//  @return map[string]interface{}
//
func (b *ByDB) GetConfig(name string) map[string]interface{} {
	cfg := config{
		Name: name,
	}
	rs := b.db.Table(b.name).Find(&cfg)
	if rs.Error == nil {
		value := map[string]interface{}{}
		_ = json.Unmarshal([]byte(cfg.Value), &value)
		return value
	}
	return nil
}

//
// SetConfig
//  @Description: 保存配置
//  @param config
//  @return bool
//  @return error
//
func (b *ByDB) SetConfig(name string, value map[string]interface{}) (bool, error) {
	cfg := config{
		Name: name,
	}
	j, _ := json.Marshal(value)
	cfg.Value = string(j)
	rs := b.db.Table(b.name).Save(cfg)
	return rs.RowsAffected > 0, rs.Error
}
