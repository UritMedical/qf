package user

import (
	"qf"
)

type GroupDal struct {
	qf.BaseDal
}

func (g GroupDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (g GroupDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//==========分组相关接口==================

//Insert 创建分组
func (g GroupDal) Insert(group *Group) error {
	return g.DB().Create(group).Error
}

//Remove 删除
func (g GroupDal) Remove(id uint) error {
	return g.DB().Where("id = ?", id).Update("remove", 1).Error
}

//Update 更新分组
func (g GroupDal) Update(id uint, groupName string) error {
	return g.DB().Where("id =?", id).Update("name", groupName).Error
}

//GetAll 获取分组列表
func (g GroupDal) GetAll() ([]Group, error) {
	list := make([]Group, 0)
	err := g.DB().Find(&list).Error
	return list, err
}

//GetGroupsByIds 根据组id列表获取组
func (g GroupDal) GetGroupsByIds(gIds []uint) ([]Group, error) {
	list := make([]Group, 0)
	err := g.DB().Where("id IN ?", gIds).Find(&list).Error
	return list, err
}

//==========组用户关系接口==================

type RelationDal struct {
	qf.BaseDal
}

func (r RelationDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (r RelationDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// AddUsersToGroup
//  @Description: 向指定组批量插入用户
//  @param gId 组Id
//  @param uIds 用户Id列表
//  @return error
//
func (r RelationDal) AddUsersToGroup(gId uint, uIds []uint) error {
	list := make([]Relation, 0)
	for _, id := range uIds {
		relation := Relation{
			GroupId: gId,
			UserId:  id,
		}
		list = append(list, relation)
	}
	return r.DB().Create(&list).Error
}

//
// RemoveUsersFromGroup
//  @Description: 从指定组移除批量用户
//  @param gId
//  @param uIds
//  @return error
//
func (r RelationDal) RemoveUsersFromGroup(gId uint, uIds []uint) error {
	return r.DB().Where("group_id = ? AND user_id IN ?", gId, uIds).Delete(&Relation{}).Error
}

//
// GetUsersByGroup
//  @Description: 获取指定组的所有用户
//  @param gId 组Id
//  @return []uint 用户Id
//  @return error
//
func (r RelationDal) GetUsersByGroup(gId uint) ([]uint, error) {
	uIds := make([]uint, 0)
	err := r.DB().Select("user_id").Where("group_id = ?", gId).Find(&uIds).Error
	return uIds, err
}

//
// GetGroupsByUser
//  @Description: 获取指定用户所在的组
//  @param uId
//  @return []uint
//  @return error
//
func (r RelationDal) GetGroupsByUser(uId uint) ([]uint, error) {
	gIds := make([]uint, 0)
	err := r.DB().Select("group_id").Where("user_id = ?", uId).Find(&gIds).Error
	return gIds, err
}
