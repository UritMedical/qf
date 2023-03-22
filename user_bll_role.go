package qf

import (
	"github.com/UritMedical/qf/util"
)

func (b *userBll) regRoleApi(api ApiMap) {
	//角色
	api.Reg(EApiKindSave, "role", b.saveRole)        //创建、修改角色
	api.Reg(EApiKindDelete, "role", b.deleteRole)    //删除角色
	api.Reg(EApiKindGetList, "roles", b.getAllRoles) //获取所有角色

	//用户-角色
	api.Reg(EApiKindSave, "role/users", b.setUserRoleRelation) //给角色删除或者添加用户
	api.Reg(EApiKindGetList, "role/users", b.getRoleUsers)     //获取指定角色下的用户

	//角色-权限组
	api.Reg(EApiKindSave, "role/permission", b.setRolePermission)      //给角色配置权限
	api.Reg(EApiKindGetList, "role/permissions", b.getRolePermissions) //获取角色拥有的权限

	//权限组
	api.Reg(EApiKindSave, "permission", b.savePermission)     //添加权限组
	api.Reg(EApiKindDelete, "permission", b.deletePermission) //删除权限组
	api.Reg(EApiKindGetList, "permissions", b.getPermissions) //获取权限组

	//权限组-API
	api.Reg(EApiKindSave, "permission/apis", b.setPermissionApi)
	api.Reg(EApiKindGetList, "permission/apis", b.getPermissionApi)
}

//
// saveRole
//  @Description: 增改角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) saveRole(ctx *Context) (interface{}, IError) {
	role := &Role{}
	if err := ctx.Bind(role); err != nil {
		return nil, err
	}
	if role.Id == 0 {
		role.Id = ctx.NewId(role)
	}
	return nil, b.roleDal.Save(role)
}

//
// deleteRole
//  @Description: 删除角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) deleteRole(ctx *Context) (interface{}, IError) {
	uId := ctx.GetId()
	return nil, b.roleDal.Delete(uId)
}

//
// getAllRoles
//  @Description: 获取所有的角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getAllRoles(_ *Context) (interface{}, IError) {
	roles := make([]Role, 0)
	err := b.roleDal.GetList(0, 100, &roles)

	result := make([]map[string]interface{}, 0)
	for _, role := range roles {
		//获取此角色拥有的用户
		userIds, _ := b.userRoleDal.GetUsersByRoleId(role.Id)
		users, _ := b.userDal.GetUsersByIds(userIds)

		//获取此角色拥有的权限
		permissionId, _ := b.rolePermissionDal.GetRolePermission(role.Id)
		permissions, _ := b.permissionDal.GetPermissionsByIds(permissionId)

		roleDict := make(map[string]interface{})
		roleDict["RoleInfo"] = util.ToMap(role)
		roleDict["Users"] = util.ToMaps(users)
		roleDict["Permissions"] = util.ToMaps(permissions)

		result = append(result, roleDict)
	}
	return result, err
}

//
// setUserRoleRelation
//  @Description: 设置角色-用户关系
//  @param roleId 角色ID
//  @param userId 用户Id
//  @return error
//
func (b *userBll) setUserRoleRelation(ctx *Context) (interface{}, IError) {
	var params = struct {
		RoleId  uint64
		UserIds []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.userRoleDal.SetRoleUsers(params.RoleId, params.UserIds)
}

//
// setRolePermission
//  @Description: 设置角色-权限关系
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) setRolePermission(ctx *Context) (interface{}, IError) {
	var params = struct {
		RoleId        uint64
		PermissionIds []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.rolePermissionDal.SetRolePermission(params.RoleId, params.PermissionIds)
}

//
// getRoleUsers
//  @Description: 获取此角色下的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getRoleUsers(ctx *Context) (interface{}, IError) {
	roleId := ctx.GetId()
	userIds, _ := b.userRoleDal.GetUsersByRoleId(roleId)
	users, err := b.userDal.GetUsersByIds(userIds)
	return util.ToMaps(users), err
}

//
// getRolePermissions
//  @Description: 获取此角色的权限
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getRolePermissions(ctx *Context) (interface{}, IError) {
	roleId := ctx.GetId()
	permissionId, _ := b.rolePermissionDal.GetRolePermission(roleId)
	permissions, err := b.permissionDal.GetPermissionsByIds(permissionId)
	return util.ToMaps(permissions), err
}

//
// getRolesByUserId
//  @Description: 获取用户所拥有的角色
//  @receiver b
//
func (b *userBll) getRolesByUserId(userId uint64) ([]Role, IError) {
	roleIds, _ := b.userRoleDal.GetRolesByUserId(userId)
	return b.roleDal.GetRolesByIds(roleIds)
}

func (b *userBll) savePermission(ctx *Context) (interface{}, IError) {
	permission := &Permission{}
	if err := ctx.Bind(permission); err != nil {
		return nil, err
	}
	if permission.Id == 0 {
		permission.Id = ctx.NewId(permission)
	}
	return nil, b.permissionDal.Save(permission)
}

func (b *userBll) deletePermission(ctx *Context) (interface{}, IError) {
	uId := ctx.GetId()
	return nil, b.permissionDal.Delete(uId)
}

func (b *userBll) getPermissions(_ *Context) (interface{}, IError) {
	permissions := make([]Permission, 0)
	err := b.permissionDal.GetList(0, 100, &permissions)
	return util.ToMaps(permissions), err
}

//
// setPermissionApi
//  @Description: 批量设置权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) setPermissionApi(ctx *Context) (interface{}, IError) {
	params := struct {
		PermissionId uint64
		ApiIds       []string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.permissionApiDal.SetPermissionApis(params.PermissionId, params.ApiIds)
}

//
// getPermissionApi
//  @Description: 获取指定权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *userBll) getPermissionApi(ctx *Context) (interface{}, IError) {
	permissionId := ctx.GetUIntValue("PermissionId")
	return b.permissionApiDal.GetApisByPermissionId(permissionId)
}

//
// getUserAllApis
//  @Description: 获取全部用户可以访问的Api列表
//  @param roleIds
//
func (b *userBll) getUserAllApis(roleIds ...uint64) map[string]byte {
	apis := map[string]byte{}
	for _, roleId := range roleIds {
		// 获取角色全部的权限列表
		permissionId, _ := b.rolePermissionDal.GetRolePermission(roleId)
		permissions, err := b.permissionDal.GetPermissionsByIds(permissionId)
		if err != nil {
			continue
		}

		// 再获取权限包含的所有Api
		for _, permission := range permissions {
			apiIds, err := b.permissionApiDal.GetApisByPermissionId(permission.Id)
			if err != nil {
				continue
			}
			for _, apiId := range apiIds {
				if _, ok := apis[apiId]; !ok {
					apis[apiId] = 1
				}
			}
		}
	}
	return apis
}
