package uDal

import "qf"

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
func (u *UserDal) CheckLogin(loginId, password string) bool {
	var count int64
	u.DB().Where("login_id = ? AND password = ?", loginId, password).Count(&count)
	return count > 0
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
