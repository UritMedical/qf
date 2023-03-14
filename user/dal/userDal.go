package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
)

type UserDal struct {
	qf.BaseDal
}

//
// SetPassword
//  @Description: 修改密码
//  @param id 用户Id
//  @param newPwd 新密码的MD5格式
//  @return error
//
func (u *UserDal) SetPassword(id uint64, newPwd string) error {
	return u.DB().Where("id = ?", id).Update("password", newPwd).Error
}

//
// CheckLogin
//  @Description: 登录检查
//  @param loginId
//  @param password
//  @return bool
//
func (u *UserDal) CheckLogin(loginId, password string) (model.User, bool) {
	var user model.User
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
func (u *UserDal) CheckOldPassword(id uint64, password string) bool {
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
func (u *UserDal) GetAllUsers() ([]model.User, error) {
	list := make([]model.User, 0)
	err := u.DB().Where("Id > 0").Find(&list).Error
	return list, err
}

//
// GetUsersByIds
//  @Description: 获取指定用户的信息
//  @param userIds
//  @return []uModel.User
//  @return error
//
func (u *UserDal) GetUsersByIds(userIds []uint64) ([]model.User, error) {
	list := make([]model.User, 0)
	err := u.DB().Where("Id IN (?)", userIds).Find(&list).Error
	return list, err
}
