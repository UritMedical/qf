package qf

//
// departmentDal
//  @Description:
//
type departmentDal struct {
	BaseDal
}

//
// GetDptsByIds
//  @Description: 获取部门列表
//  @param dptIds
//  @return []uDepartment
//  @return error
//
func (d departmentDal) GetDptsByIds(dptIds []uint64) ([]Department, IError) {
	list := make([]Department, 0)
	err := d.DB().Where("Id IN (?)", dptIds).Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

//---------------------------------------------------------------------------------------------------

//
// permissionDal
//  @Description:
//
type permissionDal struct {
	BaseDal
}

//
// GetPermissionsByIds
//  @Description: 获取权限组列表
//  @param ids
//  @return []Permission
//  @return error
//
func (p permissionDal) GetPermissionsByIds(ids []uint64) ([]Permission, IError) {
	list := make([]Permission, 0)
	err := p.DB().Where("Id IN (?)", ids).Find(&list).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return list, nil
}

//---------------------------------------------------------------------------------------------------

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

//---------------------------------------------------------------------------------------------------

type dptUserDal struct {
	BaseDal
}

//
// SetDptUsers
//  @Description: 向指定部门添加用户
//  @param departId
//  @param userIds
//  @return error
//
func (d dptUserDal) SetDptUsers(departId uint64, userIds []uint64) IError {
	oldUserIds, err := d.GetUsersByDptId(departId)
	if err != nil {
		return err
	}
	newUsers := diffIntSet(userIds, oldUserIds)
	removeUsers := diffIntSet(oldUserIds, userIds)

	tx := d.DB().Begin()
	//新增关系
	if len(newUsers) > 0 {
		addList := make([]DepartUser, 0)
		for _, id := range newUsers {
			addList = append(addList, DepartUser{
				DepartId: departId,
				UserId:   id,
			})
		}
		if err := tx.Create(&addList).Error; err != nil {
			tx.Rollback()
			return Error(ErrorCodeSaveFailure, err.Error())
		}
	}

	//删除关系
	if err := tx.Where("DepartId = ? and UserId IN (?)", departId, removeUsers).Delete(DepartUser{}).Error; err != nil {
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
func (d dptUserDal) RemoveUser(departId uint64, userId uint64) IError {
	err := d.DB().Where("DepartId = ? AND UserId = ?", departId, userId).Delete(&DepartUser{}).Error
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
func (d dptUserDal) GetUsersByDptId(departId uint64) ([]uint64, IError) {
	userIds := make([]uint64, 0)
	err := d.DB().Where("DepartId = ?", departId).Select("UserId").Find(&userIds).Error
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
func (d dptUserDal) GetDptsByUserId(userId uint64) ([]uint64, IError) {
	dptIds := make([]uint64, 0)
	err := d.DB().Where("UserId = ?", userId).Select("DepartId").Find(&dptIds).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return dptIds, nil
}

//---------------------------------------------------------------------------------------------------

type permissionApiDal struct {
	BaseDal
}

//
// SetPermissionApis
//  @Description: 向指定权限组添加，删除API
//  @param permissionId
//  @param apiKeys
//  @return error
//
func (r permissionApiDal) SetPermissionApis(permissionId uint64, apiList []ApiInfo) IError {
	tx := r.DB().Begin()
	//先删除此权限组所有的API
	if err := tx.Where("PermissionId = ?", permissionId).Delete(&PermissionApi{}).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeDeleteFailure, err.Error())
	}

	apis := make([]PermissionApi, 0)
	for _, api := range apiList {
		apis = append(apis, PermissionApi{
			PermissionId: permissionId,
			Group:        api.Group,
			ApiId:        api.ApiId,
		})
	}

	if err := tx.Create(&apis).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeSaveFailure, err.Error())
	}
	e := tx.Commit().Error
	if e != nil {
		return Error(ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// GetApisByPermissionId
//  @Description: 获取指定权限组的API
//  @param permissionId
//  @return []string
//  @return error
//
func (r permissionApiDal) GetApisByPermissionId(permissionId uint64) ([]PermissionApi, IError) {
	apis := make([]PermissionApi, 0)
	err := r.DB().Where("PermissionId = ?", permissionId).Find(&apis).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return apis, nil
}

//---------------------------------------------------------------------------------------------------

type rolePermissionDal struct {
	BaseDal
}

//
// SetRolePermission
//  @Description: 给指定角色分配权限。先删除roleId所有的权限，然后再重新添加
//  @param roleId
//  @param permissionIds
//  @return error
//
func (r rolePermissionDal) SetRolePermission(roleId uint64, permissionIds []uint64) IError {
	tx := r.DB().Begin()
	//先删除原来的权限
	if err := tx.Where("RoleId = ?", roleId).Delete(&RolePermission{}).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeDeleteFailure, err.Error())
	}

	//再将权限添加到数据库
	list := make([]RolePermission, 0)
	for _, id := range permissionIds {
		list = append(list, RolePermission{
			RoleId:       roleId,
			PermissionId: id,
		})
	}

	if err := tx.Create(&list).Error; err != nil {
		tx.Rollback()
		return Error(ErrorCodeSaveFailure, err.Error())
	}

	e := tx.Commit().Error
	if e != nil {
		return Error(ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// GetRolePermission
//  @Description: 获取指定角色拥有的权限
//  @param roleId
//
func (r rolePermissionDal) GetRolePermission(roleId uint64) ([]uint64, IError) {
	permissions := make([]uint64, 0)
	err := r.DB().Where("RoleId = ?", roleId).Select("PermissionId").Find(&permissions).Error
	if err != nil {
		return nil, Error(ErrorCodeRecordNotFound, err.Error())
	}
	return permissions, nil
}

//---------------------------------------------------------------------------------------------------

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

//---------------------------------------------------------------------------------------------------

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
