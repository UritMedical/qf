package uDal

import (
	"qf"
	uModel "qf/mc/user/model"
)

type UserDal struct {
	qf.BaseDal
}

func (u UserDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (u UserDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// SetPassword
//  @Description: 修改密码
//  @param id 用户Id
//  @param newPwd 新密码的MD5格式
//  @return error
//
func (u *UserDal) SetPassword(id uint, newPwd string) error {
	return u.DB().Where("id = ?", id).Update("password", newPwd).Error
}

//
// CheckLogin
//  @Description: 登录检查
//  @param loginId
//  @param password
//  @return bool
//
func (u *UserDal) CheckLogin(loginId, password string) (uModel.User, bool) {
	var user uModel.User
	u.DB().Where("LoginId = ? AND Password = ?", loginId, password).Find(&user)
	return user, user.Id > 0
}

//
// CheckOldPassword
//  @Description: 修改密码时旧密码检查
//  @param id
//  @param password
//  @return bool
//
func (u UserDal) CheckOldPassword(id uint, password string) bool {
	var count int64
	u.DB().Where("id = ? AND password = ?", id, password).Count(&count)
	return count > 0
}

//
// GetAllUsers
//  @Description: 获取所有用
//  @return []uModel.User
//  @return error
//
func (u UserDal) GetAllUsers() ([]uModel.User, error) {
	list := make([]uModel.User, 0)
	err := u.DB().Where("Id > 0").Find(&list).Error
	return list, err
}
