package setting

import (
	"gorm.io/gorm"
	"time"
)

func NewSettingByDB(db *gorm.DB) *ByDB {
	s := &ByDB{
		db: db,
	}
	// 生成表
	db.Table("Setting").AutoMigrate(data{})
	return s
}

type ByDB struct {
	db *gorm.DB
}

type data struct {
	CreatedAt time.Time
	Id        string
	Value     string
}

func (s *ByDB) Get(id string) string {
	data := data{}
	result := s.db.Table("Setting").Where("id = ?", id).Find(&data)
	if result.Error == nil {
		return data.Value
	}
	return ""
}

func (s *ByDB) Set(id string, value string) (bool, error) {
	data := data{
		CreatedAt: time.Now().Local(),
		Id:        id,
		Value:     value,
	}
	result := s.db.Table("Setting").Create(&data)
	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, result.Error
}
