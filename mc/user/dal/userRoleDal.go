package uDal

import (
	"qf"
	uModel "qf/mc/user/model"
	uUtils "qf/mc/user/utils"
)

type UserRoleDal struct {
	qf.BaseDal
}

func (u UserRoleDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (u UserRoleDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// SetRoleUsers
//  @Description: 增、删用户-角色关系
//  @param roleId
//  @param userIds
//  @return error
//
func (u UserRoleDal) SetRoleUsers(roleId uint, userIds []uint) error {
	oldUsers, err := u.GetUsersByRoleId(roleId)
	if err != nil {
		return err
	}
	newUsers := uUtils.DiffIntSet(userIds, oldUsers)
	removeUsers := uUtils.DiffIntSet(oldUsers, userIds)

	tx := u.DB().Begin()
	//新增关系
	if len(newUsers) > 0 {
		addList := make([]uModel.UserRole, 0)
		for _, id := range newUsers {
			addList = append(addList, uModel.UserRole{
				RoleId: roleId,
				UserId: id,
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	//删除关系
	if err := tx.Where("RoleId = ? and UserId IN (?)", roleId, removeUsers).Delete(uModel.UserRole{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
// GetUsersByRoleId
//  @Description: 获取此角色下的所有用户
//  @param roleId
//  @return []uint
//  @return error
//
func (u UserRoleDal) GetUsersByRoleId(roleId uint) ([]uint, error) {
	userIds := make([]uint, 0)
	err := u.DB().Debug().Where("RoleId = ?", roleId).Select("UserId").Find(&userIds).Error
	return userIds, err
}

//
// GetRolesByUserId
//  @Description: 获取用户所拥有的角色
//  @param userId
//  @return []uint
//  @return error
//
func (u UserRoleDal) GetRolesByUserId(userId uint) ([]uint, error) {
	roleIds := make([]uint, 0)
	err := u.DB().Where("UserId = ?", userId).Select("RoleId").Find(&roleIds).Error
	return roleIds, err
}
