package uDal

import (
	"qf"
	uModel "qf/mc/user/model"
)

type DptUserDal struct {
	qf.BaseDal
}

func (d DptUserDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (d DptUserDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// AddUsers
//  @Description: 向指定部门添加用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) AddUsers(departId uint, userIds []uint) error {
	list := make([]uModel.DepartUser, 0)
	for _, id := range userIds {
		list = append(list, uModel.DepartUser{
			DepartId: departId,
			UserId:   id,
		})
	}
	return d.DB().Create(list).Error
}

//
// RemoveUsers
//  @Description: 从指定部门移除用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) RemoveUsers(departId uint, userIds []uint) error {
	return d.DB().Where("DepartId = ? AND UserId IN (?)", departId, userIds).Delete(&uModel.DepartUser{}).Error
}

//
// GetDptUsers
//  @Description: 获取部门中所有用户
//  @param departId
//  @return []uint
//  @return error
//
func (d DptUserDal) GetDptUsers(departId uint) ([]uint, error) {
	userIds := make([]uint, 0)
	err := d.DB().Where("DepartId = ?", departId).Find(&userIds).Error
	return userIds, err
}

//
// GetUserDpts
//  @Description: 获取用户所属部门
//  @param userId
//  @return []uint
//  @return error
//
func (d DptUserDal) GetUserDpts(userId uint) ([]uint, error) {
	dptIds := make([]uint, 0)
	err := d.DB().Where("UserId = ?", userId).Find(&dptIds).Error
	return dptIds, err
}
