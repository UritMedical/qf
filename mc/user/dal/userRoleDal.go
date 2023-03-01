package uDal

import "qf"

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
//  @Description: 向指定角色添加用户
//  @param roleId 角色ID
//  @param userId 用户Id
//  @return error
//
func (u UserRoleDal) SetRoleUsers(roleId uint, userId []uint) error {
	return nil
}

//
// RemoveRole
//  @Description: 删除指定用户的角色
//  @param roleId
//  @param userId
//  @return error
//
func (u UserRoleDal) RemoveUserFromRole(roleId, userId uint) error {
	return nil
}

//
// GetUserRole
//  @Description: 获取用户拥有的角色
//  @param userId
//  @return []uint
//  @return error
//
func (u UserRoleDal) GetUserRole(userId uint) ([]uint, error) {
	return nil, nil
}
