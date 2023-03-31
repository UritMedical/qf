package qid

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

//
// NewIdAllocatorByDB
//  @Description:
//  @param per
//  @param path
//  @return *ByDB
//
func NewIdAllocatorByDB(per uint, start uint, db *gorm.DB) IIdAllocator {
	name := "QfId"
	if db.Migrator().HasTable(name) == false {
		_ = db.Table(name).AutoMigrate(idAllocator{})
	}
	return &byDB{
		name:  name,
		db:    db,
		per:   per,
		start: start,
		lock:  &sync.RWMutex{},
	}
}

type byDB struct {
	name  string
	db    *gorm.DB
	per   uint
	start uint
	lock  *sync.RWMutex
}

type idAllocator struct {
	Name     string `gorm:"primaryKey"`
	Value    uint64
	LastTime time.Time
}

//
// Next
//  @Description: 下一个Id
//  @param name 名称
//  @return uint64
//
func (b *byDB) Next(name string) uint64 {
	b.lock.Lock()
	defer b.lock.Unlock()

	// 先获取
	val := idAllocator{
		Name: name,
	}
	rs := b.db.Table(b.name).Find(&val)
	if rs.RowsAffected == 0 {
		val.Value = uint64(b.start)
	} else {
		val.Value += 1
	}
	// 再保存
	val.Name = name
	val.LastTime = time.Now().Local()
	rs = b.db.Table(b.name).Save(&val)
	if rs.RowsAffected > 0 {
		if b.per == 0 {
			return val.Value
		}
		// 附加前缀
		nid, _ := strconv.Atoi(fmt.Sprintf("%d%d", b.per, val.Value))
		return uint64(nid)
	}
	return 0
}
