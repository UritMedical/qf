package patient

import (
	"gorm.io/gorm"
	"qf/helper/content"
	"time"
)

type dal struct {
	dbNoIndex   *gorm.DB
	dbNameIndex *gorm.DB
}

type noIndex struct {
	PatNo     string `gorm:"primarykey"`
	ContentId uint   `gorm:"primarykey"`
	Time      time.Time
}

type nameIndex struct {
	PatName   string `gorm:"primarykey"`
	ContentId uint   `gorm:"primarykey"`
	Time      time.Time
}

func newDal(db *gorm.DB) *dal {
	d := &dal{}
	// 创建索引表
	d.dbNoIndex = db.Table("patient_index_no")
	_ = d.dbNoIndex.AutoMigrate(noIndex{})
	d.dbNameIndex = db.Table("patient_index_name")
	_ = d.dbNameIndex.AutoMigrate(nameIndex{})
	return d
}

// 更新索引
//  old: 旧内容
//  latest: 新内容
func (d *dal) updateIndexes(old content.Content, latest content.Content) error {
	// 获取最新提交的内容
	var pat struct {
		No   string // 患者病历号（门诊号/住院号/体检号/...）
		Name string // 姓名
	}
	err := latest.BindJson(&pat)
	if err != nil {
		return err
	}
	// 创建病历号索引
	err = d.dbNoIndex.Transaction(func(tx *gorm.DB) error {
		// 先删除旧记录记录
		rs := tx.Delete(noIndex{
			PatNo:     pat.No,
			Time:      old.Time,
			ContentId: old.ID,
		})
		if rs.Error != nil {
			return rs.Error
		}
		// 再写入新记录
		rs = tx.Create(noIndex{
			PatNo:     pat.No,
			Time:      latest.Time,
			ContentId: latest.ID,
		})
		if rs.Error != nil {
			return rs.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 创建姓名索引
	err = d.dbNameIndex.Transaction(func(tx *gorm.DB) error {
		// 先删除旧记录记录
		rs := tx.Delete(nameIndex{
			PatName:   pat.No,
			Time:      old.Time,
			ContentId: old.ID,
		})
		if rs.Error != nil {
			return rs.Error
		}
		// 再写入新记录
		rs = tx.Create(nameIndex{
			PatName:   pat.No,
			Time:      latest.Time,
			ContentId: latest.ID,
		})
		if rs.Error != nil {
			return rs.Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (d *dal) selectList(query map[string]interface{}) (interface{}, error) {
	return nil, nil
}
