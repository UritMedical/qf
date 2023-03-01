package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (u *UserBll) regRoleApi(api qf.ApiMap) {
	//角色
	api.Reg(qf.EKindSave, "role", u.saveRole)        //创建、修改角色
	api.Reg(qf.EKindDelete, "role", u.deleteRole)    //删除角色
	api.Reg(qf.EKindGetList, "roles", u.getAllRoles) //获取所有就是

	//用户-角色
	api.Reg(qf.EKindSave, "role/users", u.setRoleUsers)        //给角色设置用户，删除或者添加
	api.Reg(qf.EKindDelete, "role/user", u.removeUserFromRole) //从指定角色中删除用户

	//角色-权限组
	api.Reg(qf.EKindSave, "role/rights", u.setRoleRights)    //给角色配置权限
	api.Reg(qf.EKindGetList, "role/rights", u.getRoleRights) //获取角色拥有的权限
}

func (u *UserBll) saveRole(ctx *qf.Context) (interface{}, error) {
	var role uModel.Role
	if err := ctx.Bind(&role); err != nil {
		return nil, err
	}
	return nil, u.roleDal.Save(&role)
}

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
// SetRoleUsers
//  @Description: 向指定角色添加用户
//  @param roleId 角色ID
//  @param userId 用户Id
//  @return error
//
func (u *UserBll) setRoleUsers(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId  uint
		UserIds []uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.userRole.SetRoleUsers(params.RoleId, params.UserIds)
}

func (u *UserBll) removeUserFromRole(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId uint
		UserId uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.userRole.RemoveUserFromRole(params.RoleId, params.UserId)
}

func (u *UserBll) setRoleRights(ctx *qf.Context) (interface{}, error) {
	var params = struct {
		RoleId    uint
		RightsIds []uint
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.roleRightsDal.SetRoleRights(params.RoleId, params.RightsIds)
}

func (u *UserBll) getRoleRights(ctx *qf.Context) (interface{}, error) {
	roleId := ctx.GetUIntValue("RoleId")
	return u.roleRightsDal.GetRoleRights(roleId)
}
