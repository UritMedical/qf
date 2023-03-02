package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (u *UserBll) regRoleApi(api qf.ApiMap) {
	//角色
	api.Reg(qf.EKindSave, "role", u.saveRole)        //创建、修改角色
	api.Reg(qf.EKindDelete, "role", u.deleteRole)    //删除角色
	api.Reg(qf.EKindGetList, "roles", u.getAllRoles) //获取所有角色

	//用户-角色
	api.Reg(qf.EKindSave, "role/users", u.setUserRoleRelation) //给角色删除或者添加用户
	api.Reg(qf.EKindGetList, "role/users", u.getRoleUsers)     //获取指定角色下的用户

	//角色-权限组
	api.Reg(qf.EKindSave, "role/rights", u.setRoleRightsRelation) //给角色配置权限
	api.Reg(qf.EKindGetList, "role/rights/", u.getRoleRights)     //获取角色拥有的权限
}

//
// saveRole
//  @Description: 增改角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) saveRole(ctx *qf.Context) (interface{}, error) {
	role := &uModel.Role{}
	if err := ctx.Bind(role); err != nil {
		return nil, err
	}
	return nil, u.roleDal.Save(role)
}

//
// deleteRole
//  @Description: 删除角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) deleteRole(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.roleDal.Delete(uId)
	return nil, err
}

//
// getAllRoles
//  @Description: 获取所有的角色
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getAllRoles(ctx *qf.Context) (interface{}, error) {
	roles := make([]uModel.Role, 0)
	err := u.roleDal.GetList(0, 100, &roles)
	return roles, err
}

//
// setUserRoleRelation
//  @Description: 设置角色-用户关系
//  @param roleId 角色ID
//  @param userId 用户Id
//  @return error
//
func (u *UserBll) setUserRoleRelation(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId  uint
		UserIds []uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.userRoleDal.SetRoleUsers(params.RoleId, params.UserIds)
}

//
// setRoleRightsRelation
//  @Description: 设置角色-权限关系
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) setRoleRightsRelation(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId    uint
		RightsIds []uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.roleRightsDal.SetRoleRights(params.RoleId, params.RightsIds)
}

//
// getRoleUsers
//  @Description: 获取此角色下的用户
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getRoleUsers(ctx *qf.Context) (interface{}, error) {
	roleId := ctx.GetUIntValue("RoleId")
	userIds, err := u.userRoleDal.GetUsersByRoleId(uint(roleId))
	//TODO 把用户Id转换成用户信息
	return userIds, err
}

//
// getRoleRights
//  @Description: 获取此角色的权限
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *UserBll) getRoleRights(ctx *qf.Context) (interface{}, error) {
	roleId := ctx.GetUIntValue("RoleId")
	rightsId, err := u.roleRightsDal.GetRoleRights(uint(roleId))
	//TODO 把角色Id转换成角色信息
	return rightsId, err
}
