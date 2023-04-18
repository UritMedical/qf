package qf

import "strings"

//region base dict dal

//
// departmentDal
//  @Description:
//
type departmentDal struct {
	BaseDal
}

//
// GetAll
//  @Description: 获取所有部门列表
//  @return []Dept
//  @return IError
//
func (d departmentDal) GetAll() ([]Dept, IError) {
	list := make([]Dept, 0)
	err := d.DB().Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

// roleDal
//  @Description: 角色
//
type roleDal struct {
	BaseDal
}

//
// GetRolesByIds
//  @Description: 获取角色列表
//  @param ids
//  @return []uRole
//  @return error
//
func (role roleDal) GetRolesByIds(ids []uint64) ([]Role, IError) {
	list := make([]Role, 0)
	err := role.DB().Where("Id IN (?)", ids).Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

//endregion

//region role api dal

type roleApiDal struct {
	BaseDal
}

//
// SetRoleApis
//  @Description: 向指定角色添加，删除API
//  @param roleId
//  @param apiKeys
//  @return error
//
func (r roleApiDal) SetRoleApis(roleId uint64, apiList []string) IError {
	tx := r.DB().Begin()
	//先删除此角色组所有的API
	if err := tx.Where("RoleId = ?", roleId).Delete(&RoleApi{}).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeDeleteFailure, err.Error())
	}

	if len(apiList) > 0 {
		addList := make([]RoleApi, 0)
		for _, v := range apiList {
			ret := strings.Split(v, ":")
			addList = append(addList, RoleApi{
				RoleId:     roleId,
				Permission: ret[0],
				Url:        ret[1],
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return Error(ErrorCodeSaveFailure, err.Error())
		}
	}
	e := tx.Commit().Error
	if e != nil {
		return Error(ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// GetApisByRoleId
//  @Description: 获取指定角色的API
//  @param roleId
//  @return []string
//  @return error
//
func (r roleApiDal) GetApisByRoleId(roleIds []uint64) ([]RoleApi, IError) {
	apis := make([]RoleApi, 0)
	err := r.DB().Where("RoleId IN ?", roleIds).Find(&apis).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return apis, nil
}

//endregion

//region user dal

type userDal struct {
	BaseDal
}

//
// SetPassword
//  @Description: 修改密码
//  @param id 用户Id
//  @param newPwd 新密码的MD5格式
//  @return error
//
func (u *userDal) SetPassword(id uint64, newPwd string) IError {
	err := u.DB().Where("id = ?", id).Update("password", newPwd).Error
	if err != nil {
		return Error(ErrorCodeSaveFailure, err.Error())
	}
	return nil
}

//
// CheckLogin
//  @Description: 登录检查
//  @param loginId
//  @param password
//  @return bool
//
func (u *userDal) CheckLogin(loginId, password string) (User, bool) {
	var user User
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
func (u *userDal) CheckOldPassword(id uint64, password string) bool {
	var count int64
	u.DB().Where("id = ? AND password = ?", id, password).Count(&count)
	return count > 0
}

//
// GetAllUsers
//  @Description: 获取所有用，不返回admin账号
//  @return []uUser
//  @return error
//
func (u *userDal) GetAllUsers() ([]User, IError) {
	list := make([]User, 0)
	err := u.DB().Where("Id > 2").Order("Id ASC").Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

//
// GetUsersByIds
//  @Description: 获取指定用户的信息
//  @param userIds
//  @return []uUser
//  @return error
//
func (u *userDal) GetUsersByIds(userIds []uint64) ([]User, IError) {
	list := make([]User, 0)
	err := u.DB().Where("Id IN (?)", userIds).Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

//endregion

//region user role dal

type userRoleDal struct {
	BaseDal
}

//
// SetRoleUsers
//  @Description: 增、删用户-角色关系
//  @param roleId
//  @param userIds
//  @return error
//
func (u *userRoleDal) SetRoleUsers(roleId uint64, userIds []uint64) IError {
	oldUsers, err := u.GetUsersByRoleId(roleId)
	if err != nil {
		return err
	}
	newUsers := diffIntSet(userIds, oldUsers)
	removeUsers := diffIntSet(oldUsers, userIds)

	tx := u.DB().Begin()
	//新增关系
	if len(newUsers) > 0 {
		addList := make([]UserRole, 0)
		for _, id := range newUsers {
			addList = append(addList, UserRole{
				RoleId: roleId,
				UserId: id,
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return Error(ErrorCodeSaveFailure, err.Error())
		}
	}

	//删除关系
	if err := tx.Where("RoleId = ? and UserId IN (?)", roleId, removeUsers).Delete(UserRole{}).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeDeleteFailure, err.Error())
	}
	e := tx.Commit().Error
	if e != nil {
		return Error(ErrorCodeSaveFailure, e.Error())
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
func (u *userRoleDal) GetUsersByRoleId(roleId uint64) ([]uint64, IError) {
	userIds := make([]uint64, 0)
	err := u.DB().Debug().Where("RoleId = ?", roleId).Select("UserId").Find(&userIds).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
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
func (u *userRoleDal) GetRolesByUserId(userId uint64) ([]uint64, IError) {
	roleIds := make([]uint64, 0)
	err := u.DB().Where("UserId = ?", userId).Select("RoleId").Find(&roleIds).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return roleIds, nil
}

//endregion

//region user dp dal

type userDpDal struct {
	BaseDal
}

//
// SetDptUsers
//  @Description: 向指定部门添加用户
//  @param departId
//  @param userIds
//  @return error
//
func (d userDpDal) SetDptUsers(departId uint64, userIds []uint64) IError {
	oldUserIds, err := d.GetUsersByDptId(departId)
	if err != nil {
		return err
	}
	newUsers := diffIntSet(userIds, oldUserIds)
	removeUsers := diffIntSet(oldUserIds, userIds)

	tx := d.DB().Begin()
	//新增关系
	if len(newUsers) > 0 {
		addList := make([]UserDept, 0)
		for _, id := range newUsers {
			addList = append(addList, UserDept{
				DeptId: departId,
				UserId: id,
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return Error(ErrorCodeSaveFailure, err.Error())
		}
	}

	//删除关系
	if err := tx.Where("DeptId = ? and UserId IN (?)", departId, removeUsers).Delete(UserDept{}).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeDeleteFailure, err.Error())
	}
	e := tx.Commit().Error
	if e != nil {
		return Error(ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// RemoveUser
//  @Description: 从指定部门移除用户
//  @param departId
//  @param userIds
//  @return error
//
func (d userDpDal) RemoveUser(departId uint64, userId uint64) IError {
	err := d.DB().Where("DeptId = ? AND UserId = ?", departId, userId).Delete(&UserDept{}).Error
	if err != nil {
		return Error(ErrorCodeDeleteFailure, err.Error())
	}
	return nil
}

//
// GetUsersByDptId
//  @Description: 获取部门中所有用户
//  @param departId
//  @return []uint64
//  @return error
//
func (d userDpDal) GetUsersByDptId(departId uint64) ([]uint64, IError) {
	userIds := make([]uint64, 0)
	err := d.DB().Where("DeptId = ?", departId).Select("UserId").Find(&userIds).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return userIds, nil
}

//
// GetDptsByUserId
//  @Description: 获取用户所属部门
//  @param userId
//  @return []uint64
//  @return error
//
func (d userDpDal) GetDptsByUserId(userId uint64) ([]uint64, IError) {
	dptIds := make([]uint64, 0)
	err := d.DB().Where("UserId = ?", userId).Select("DeptId").Find(&dptIds).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return dptIds, nil
}

//endregion
