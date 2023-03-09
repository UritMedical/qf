package uDal

import (
	"github.com/Urit-Mediacal/qf"
	"github.com/Urit-Mediacal/qf/mc/user/uModel"
	uUtils "github.com/Urit-Mediacal/qf/mc/user/utils"
)

type DptUserDal struct {
	qf.BaseDal
}

//
// AddUsers
//  @Description: 向指定部门添加用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) AddUsers(departId uint64, userIds []uint64) error {
	oldUserIds, err := d.GetUsersByDptId(departId)
	if err != nil {
		return err
	}

	//过滤出部门中已经存在的账号
	newUsers := uUtils.DiffIntSet(userIds, oldUserIds)

	list := make([]uModel.DepartUser, 0)
	for _, id := range newUsers {
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
func (d DptUserDal) RemoveUser(departId uint64, userId uint64) error {
	return d.DB().Where("DepartId = ? AND UserId = ?", departId, userId).Delete(&uModel.DepartUser{}).Error
}

//
// GetUsersByDptId
//  @Description: 获取部门中所有用户
//  @param departId
//  @return []uint64
//  @return error
//
func (d DptUserDal) GetUsersByDptId(departId uint64) ([]uint64, error) {
	userIds := make([]uint64, 0)
	err := d.DB().Where("DepartId = ?", departId).Select("UserId").Find(&userIds).Error
	return userIds, err
}

//
// GetDptsByUserId
//  @Description: 获取用户所属部门
//  @param userId
//  @return []uint64
//  @return error
//
func (d DptUserDal) GetDptsByUserId(userId uint64) ([]uint64, error) {
	dptIds := make([]uint64, 0)
	err := d.DB().Where("UserId = ?", userId).Select("DepartId").Find(&dptIds).Error
	return dptIds, err
}
