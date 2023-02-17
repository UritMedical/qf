package content

import (
	"gorm.io/gorm"
	"time"
)

func NewContentByDB(id string, db *gorm.DB) *ByDB {
	_ = db.Table(id).AutoMigrate(data{})
	return &ByDB{id: id, db: db}
}

type ByDB struct {
	id string
	db *gorm.DB
}

type data struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      string
	Info      string
}

func (c *ByDB) Insert(cnt Content) (Content, error) {
	data := data{
		CreatedAt: cnt.Time,
		User:      cnt.User,
		Info:      cnt.Info,
	}
	result := c.db.Table(c.id).Create(&data)
	if result.RowsAffected > 0 {
		cnt.ID = data.ID
		return cnt, nil
	}
	return Content{}, result.Error
}

func (c *ByDB) Update(cnt Content) (Content, error) {
	nd := data{
		CreatedAt: cnt.Time,
		User:      cnt.User,
		Info:      cnt.Info,
	}
	err := c.db.Transaction(func(tx *gorm.DB) error {
		// 先删除当前记录
		del := data{
			ID:        cnt.ID,
			DeletedAt: gorm.DeletedAt{},
		}
		del.DeletedAt.Time = cnt.Time
		rs := tx.Model(&data{}).Delete(del)
		if rs.Error != nil {
			return rs.Error
		}
		// 在写入新记录
		rs = c.db.Create(&nd)
		if rs.Error != nil {
			return rs.Error
		}
		return nil
	})
	cnt.ID = nd.ID
	return cnt, err
}

func (c *ByDB) Save(cnt Content) (Content, error) {
	m, err := c.GetModel(cnt.ID)
	if err != nil {
		return Content{}, err
	}
	if m.ID == 0 {
		return c.Insert(cnt)
	}
	return c.Update(cnt)
}

func (c *ByDB) Delete(id uint) error {
	model := data{ID: id}
	model.DeletedAt.Time = time.Now().Local()
	result := c.db.Table(c.id).Model(&data{}).Where("id = ?", model.ID).Updates(model)
	if result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}

func (c *ByDB) GetModel(id uint) (Content, error) {
	model := data{}
	result := c.db.Table(c.id).Where("id = ?", id).Find(&model)
	// 如果异常或者未查询到任何数据
	if result.RowsAffected == 0 || result.Error != nil {
		return Content{
			ID:   model.ID,
			Time: model.CreatedAt,
			User: model.User,
			Info: model.Info,
		}, result.Error
	}
	return Content{}, nil
}

func (c *ByDB) GetList(startTime, endTime time.Time) ([]Content, error) {
	list := make([]data, 0)
	result := c.db.Table(c.id).Where("created_at >= ? and created_at <= ?", startTime, endTime).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	finals := make([]Content, 0)
	for _, model := range list {
		finals = append(finals, Content{
			ID:   model.ID,
			Time: model.CreatedAt,
			User: model.User,
			Info: model.Info,
		})
	}
	return finals, nil
}
