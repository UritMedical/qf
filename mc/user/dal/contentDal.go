package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/user/model"
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
func (d DepartmentDal) GetDptsByIds(dptIds []uint64) ([]model.Department, error) {
	list := make([]model.Department, 0)
	err := d.DB().Where("Id IN (?)", dptIds).Find(&list).Error
	return list, err
}

// PermissionDal
// @Description: 权限组
//
type PermissionDal struct {
	qf.BaseDal
}

//
// GetPermissionsByIds
//  @Description: 获取权限组列表
//  @param ids
//  @return []model.PermissionGroup
//  @return error
//
func (p PermissionDal) GetPermissionsByIds(ids []uint64) ([]model.PermissionGroup, error) {
	list := make([]model.PermissionGroup, 0)
	err := p.DB().Where("Id IN (?)", ids).Find(&list).Error
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
func (role RoleDal) GetRolesByIds(ids []uint64) ([]model.Role, error) {
	list := make([]model.Role, 0)
	err := role.DB().Where("Id IN (?)", ids).Find(&list).Error
	return list, err
}
