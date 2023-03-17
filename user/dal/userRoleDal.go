package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
	"github.com/UritMedical/qf/util"
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
func (u *UserRoleDal) SetRoleUsers(roleId uint64, userIds []uint64) qf.IError {
	oldUsers, err := u.GetUsersByRoleId(roleId)
	if err != nil {
		return err
	}
	newUsers := util.DiffIntSet(userIds, oldUsers)
	removeUsers := util.DiffIntSet(oldUsers, userIds)

	tx := u.DB().Begin()
	//新增关系
	if len(newUsers) > 0 {
		addList := make([]model.UserRole, 0)
		for _, id := range newUsers {
			addList = append(addList, model.UserRole{
				RoleId: roleId,
				UserId: id,
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return qf.Error(qf.ErrorCodeSaveFailure, err.Error())
		}
	}

	//删除关系
	if err := tx.Where("RoleId = ? and UserId IN (?)", roleId, removeUsers).Delete(model.UserRole{}).Error; err != nil {
		tx.Rollback()
		return qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
	}
	e := tx.Commit().Error
	if e != nil {
		return qf.Error(qf.ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// GetUsersByRoleId
//  @Description: 获取此角色下的所有用户
//  @param roleId
//  @return []uint64
//  @return error
//
func (u *UserRoleDal) GetUsersByRoleId(roleId uint64) ([]uint64, qf.IError) {
	userIds := make([]uint64, 0)
	err := u.DB().Debug().Where("RoleId = ?", roleId).Select("UserId").Find(&userIds).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return userIds, nil
}

//
// GetRolesByUserId
//  @Description: 获取用户所拥有的角色
//  @param userId
//  @return []uint64
//  @return error
//
func (u *UserRoleDal) GetRolesByUserId(userId uint64) ([]uint64, qf.IError) {
	roleIds := make([]uint64, 0)
	err := u.DB().Where("UserId = ?", userId).Select("RoleId").Find(&roleIds).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return roleIds, nil
}
