package uDal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/user/uModel"
	uUtils "github.com/UritMedical/qf/mc/user/utils"
)

type UserRoleDal struct {
	qf.BaseDal
}

//
// SetRoleUsers
//  @Description: 增、删用户-角色关系
//  @param roleId
//  @param userIds
//  @return error
//
func (u *UserRoleDal) SetRoleUsers(roleId uint64, userIds []uint64) error {
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
//  @return []uint64
//  @return error
//
func (u *UserRoleDal) GetUsersByRoleId(roleId uint64) ([]uint64, error) {
	userIds := make([]uint64, 0)
	err := u.DB().Debug().Where("RoleId = ?", roleId).Select("UserId").Find(&userIds).Error
	return userIds, err
}

//
// GetRolesByUserId
//  @Description: 获取用户所拥有的角色
//  @param userId
//  @return []uint64
//  @return error
//
func (u *UserRoleDal) GetRolesByUserId(userId uint64) ([]uint64, error) {
	roleIds := make([]uint64, 0)
	err := u.DB().Where("UserId = ?", userId).Select("RoleId").Find(&roleIds).Error
	return roleIds, err
}
