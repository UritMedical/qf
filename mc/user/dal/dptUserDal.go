package uDal

import "qf"

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
	return nil
}

//
// RemoveUsers
//  @Description: 从指定部门移除用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) RemoveUsers(departId uint, userIds []uint) error {
	return nil
}
