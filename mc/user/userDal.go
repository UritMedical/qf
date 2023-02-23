package user

import (
	"qf"
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
// Insert
//  @Description: 创建账号
//  @param user 用户结构体
//  @return error
//
func (u *UserDal) Insert(user *User) error {
	return u.DB().Create(user).Error
}

//
// ChangePassword
//  @Description: 修改密码
//  @param id 用户Id
//  @param newPwd 新密码的MD5格式
//  @return error
//
func (u *UserDal) ChangePassword(id uint, newPwd string) error {
	return u.DB().Where("id = ?", id).Update("password", newPwd).Error
}

//
// Remove
//  @Description: 删除账号
//  @param id 用户ID
//  @return error
//
func (u *UserDal) Remove(id uint) error {
	return u.DB().Where("id = ?", id).Update("remove", 1).Error
}

//
// SetRole
//  @Description: 设置权限
//  @param id 用Id
//  @param role 角色
//  @return error
//
func (u *UserDal) SetRole(id uint, role byte) error {
	return u.DB().Where("id = ?", id).Update("role", role).Error
}

//
// ChangeAccountInfo
//  @Description: 修改登录账号
//  @param id 用户Id
//  @param newLoginId 新的登录Id
//  @param newName 新的姓名
//  @return error
//
func (u *UserDal) ChangeAccountInfo(id uint, newLoginId, newName string) error {
	return u.DB().Where("id = ?", id).Updates(map[string]interface{}{
		"login_id": newLoginId,
		"name":     newName,
	}).Error
}

//
// UpdateMisc
//  @Description: 更新其他信息
//  @param id 用户ID
//  @param misc 信息对象的Json格式
//  @return error
//
func (u *UserDal) UpdateMisc(id uint, misc string) error {
	return u.DB().Where("id = ?", id).Update("misc", misc).Error
}

//
// GetUser
//  @Description: 根据用户Id获取用户信息
//  @param id 用户Id
//  @return *User 用户信息
//  @return error
//
func (u *UserDal) GetUser(id uint) (*User, error) {
	var user User
	err := u.DB().Where("id = ?", id).Find(&user).Error
	return &user, err
}

//
// GetUserByLoginId
//  @Description: 根据账号获取用户信息
//  @param loginId 登录Id
//  @return *User 用户信息
//  @return error
//
func (u *UserDal) GetUserByLoginId(loginId string) (*User, error) {
	var user User
	err := u.DB().Where("login_id = ?", loginId).Find(&user).Error
	return &user, err
}

//
// GetAllUser
//  @Description: 获取所有用户列表
//  @return *[]User 所有用户列表
//  @return error
//
func (u *UserDal) GetAllUser() (*[]User, error) {
	list := make([]User, 0)
	err := u.DB().Where("remove = 0").Find(&list).Error
	return &list, err
}
