package uDal

import (
	"qf"
	"qf/mc/user/uModel"
)

// DepartmentDal
// @Description: 部门
//
type DepartmentDal struct {
	qf.BaseDal
}

//
// GetDptsByIds
//  @Description: 获取部门列表
//  @param dptIds
//  @return []uModel.Department
//  @return error
//
func (d DepartmentDal) GetDptsByIds(dptIds []uint64) ([]uModel.Department, error) {
	list := make([]uModel.Department, 0)
	err := d.DB().Where("Id IN (?)", dptIds).Find(&list).Error
	return list, err
}

// RightsGroupDal
// @Description: 权限组
//
type RightsGroupDal struct {
	qf.BaseDal
}

//
// GetRightsGroupByIds
//  @Description: 获取权限组列表
//  @param ids
//  @return []uModel.RightsGroup
//  @return error
//
func (r RightsGroupDal) GetRightsGroupByIds(ids []uint64) ([]uModel.RightsGroup, error) {
	list := make([]uModel.RightsGroup, 0)
	err := r.DB().Where("Id IN (?)", ids).Find(&list).Error
	return list, err
}

// RoleDal
// @Description: 角色
//
type RoleDal struct {
	qf.BaseDal
}

//
// GetRolesByIds
//  @Description: 获取角色列表
//  @param ids
//  @return []uModel.Role
//  @return error
//
func (role RoleDal) GetRolesByIds(ids []uint64) ([]uModel.Role, error) {
	list := make([]uModel.Role, 0)
	err := role.DB().Where("Id IN (?)", ids).Find(&list).Error
	return list, err
}
