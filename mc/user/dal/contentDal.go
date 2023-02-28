package uDal

import "qf"

// DepartmentDal
// @Description: 部门
//
type DepartmentDal struct {
	qf.BaseDal
}

func (o DepartmentDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (o DepartmentDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

// RightsGroupDal
// @Description: 权限组
//
type RightsGroupDal struct {
	qf.BaseDal
}

func (r RightsGroupDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (r RightsGroupDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

// RoleDal
// @Description: 角色
//
type RoleDal struct {
	qf.BaseDal
}

func (r RoleDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (r RoleDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}
