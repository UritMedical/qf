package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (b *Bll) regRoleApi(api qf.ApiMap) {
	//角色
	api.Reg(qf.EApiKindSave, "role", b.saveRole)        //创建、修改角色
	api.Reg(qf.EApiKindDelete, "role", b.deleteRole)    //删除角色
	api.Reg(qf.EApiKindGetList, "roles", b.getAllRoles) //获取所有角色

	//用户-角色
	api.Reg(qf.EApiKindSave, "role/users", b.setUserRoleRelation) //给角色删除或者添加用户
	api.Reg(qf.EApiKindGetList, "role/users", b.getRoleUsers)     //获取指定角色下的用户

	//角色-权限组
	api.Reg(qf.EApiKindSave, "role/rights", b.setRoleRightsRelation) //给角色配置权限
	api.Reg(qf.EApiKindGetList, "role/rights", b.getRoleRights)      //获取角色拥有的权限
}

//
// saveRole
//  @Description: 增改角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) saveRole(ctx *qf.Context) (interface{}, error) {
	role := &uModel.Role{}
	if err := ctx.Bind(role); err != nil {
		return nil, err
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
func (b *Bll) deleteRole(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	ret, err := b.roleDal.Delete(uId)
	return ret, err
}

//
// getAllRoles
//  @Description: 获取所有的角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getAllRoles(ctx *qf.Context) (interface{}, error) {
	roles := make([]uModel.Role, 0)
	err := b.roleDal.GetList(0, 100, &roles)
	return b.Maps(roles), err
}

//
// setUserRoleRelation
//  @Description: 设置角色-用户关系
//  @param roleId 角色ID
//  @param userId 用户Id
//  @return error
//
func (b *Bll) setUserRoleRelation(ctx *qf.Context) (interface{}, error) {
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
// setRoleRightsRelation
//  @Description: 设置角色-权限关系
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) setRoleRightsRelation(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId    uint64
		RightsIds []uint64
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.roleRightsDal.SetRoleRights(params.RoleId, params.RightsIds)
}

//
// getRoleUsers
//  @Description: 获取此角色下的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getRoleUsers(ctx *qf.Context) (interface{}, error) {
	roleId := ctx.GetUIntValue("RoleId")
	userIds, _ := b.userRoleDal.GetUsersByRoleId(roleId)
	users, err := b.userDal.GetUsersByIds(userIds)
	return b.Maps(users), err
}

//
// getRoleRights
//  @Description: 获取此角色的权限
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getRoleRights(ctx *qf.Context) (interface{}, error) {
	roleId := ctx.GetUIntValue("RoleId")
	rightsId, _ := b.roleRightsDal.GetRoleRights(roleId)
	rights, err := b.rightsDal.GetRightsGroupByIds(rightsId)
	return b.Maps(rights), err
}
