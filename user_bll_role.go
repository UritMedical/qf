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

	//角色-API
	api.Reg(EApiKindSave, "role/apis", b.setRoleApi)    //设置角色所拥有的api
	api.Reg(EApiKindGetList, "role/apis", b.getRoleApi) //获取角色拥有的API
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

		roleDict := make(map[string]interface{})
		roleDict["RoleInfo"] = util.ToMap(role)
		roleDict["Users"] = util.ToMaps(users)

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
// getRolesByUserId
//  @Description: 获取用户所拥有的角色
//  @receiver b
//
func (b *userBll) getRolesByUserId(userId uint64) ([]Role, IError) {
	roleIds, _ := b.userRoleDal.GetRolesByUserId(userId)
	return b.roleDal.GetRolesByIds(roleIds)
}

//
// getUserAllApis
//  @Description: 获取全部用户可以访问的Api列表
//  @param roleIds
//
func (b *userBll) getUserAllApis(roles []RoleInfo) map[string]byte {
	apis := map[string]byte{}
	//for _, r := range roles {
	//	// 获取角色全部的权限列表
	//	permissionId, _ := b.rolePermissionDal.GetRolePermission(r.Id)
	//	permissions, err := b.permissionDal.GetPermissionsByIds(permissionId)
	//	if err != nil {
	//		continue
	//	}
	//
	//	// 再获取权限包含的所有Api
	//	for _, permission := range permissions {
	//		list, err := b.roleApiDal.GetApisByPermissionId(permission.Id)
	//		if err != nil {
	//			continue
	//		}
	//		for _, v := range list {
	//			if _, ok := apis[v.ApiId]; !ok {
	//				apis[v.ApiId] = 1
	//			}
	//		}
	//	}
	//}
	return apis
}

//设置角色所拥有的api
func (b *userBll) setRoleApi(ctx *Context) (interface{}, IError) {
	params := struct {
		RoleId uint64
		Apis   []string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.roleApiDal.SetRoleApis(params.RoleId, params.Apis)
}

//获取角色所拥有的api
func (b *userBll) getRoleApi(ctx *Context) (interface{}, IError) {
	roleId := ctx.GetId()
	ret, err := b.roleApiDal.GetApisByRoleId([]uint64{roleId})
	if err != nil {
		return nil, err
	}
	return util.ToMaps(ret), nil
}
