package id

import (
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type ByDB struct {
	db   *gorm.DB
	per  uint
	lock *sync.RWMutex
}

type idAllocator struct {
	Name     string `gorm:"primaryKey"`
	Value    uint64
	LastTime time.Time
}

//
// NewIdAllocatorByDB
//  @Description:
//  @param per
//  @param path
//  @return *ByDB
//
func NewIdAllocatorByDB(per uint, db *gorm.DB) *ByDB {
	_ = db.Table("id_allocator").AutoMigrate(idAllocator{})
	return &ByDB{
		db:   db,
		per:  per,
		lock: &sync.RWMutex{},
	}
}

//
// Next
//  @Description: 下一个Id
//  @receiver b
//  @param name
//  @return uint64
//
func (b *ByDB) Next(name string) uint64 {
	b.lock.Lock()
	defer b.lock.Unlock()

	// 先获取
	val := idAllocator{
		Name: name,
	}
	rs := b.db.Table("id_allocator").First(&val)
	if rs.RowsAffected == 0 {
		val.Value = 1
	} else {
		val.Value += 1
	}
	// 再保存
	val.Name = name
	val.LastTime = time.Now().Local()
	rs = b.db.Table("id_allocator").Save(&val)
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
