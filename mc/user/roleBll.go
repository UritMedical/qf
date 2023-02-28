package user

import "qf"

func (u *UserBll) regRoleApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "role", u.saveRole)
	api.Reg(qf.EKindDelete, "role", u.deleteRole)
	api.Reg(qf.EKindGetModel, "roles", u.getRightGroupList)
}

func (u *UserBll) saveRole(ctx *qf.Context) (interface{}, error) {
	var role Role
	if err := ctx.BindModel(&role); err != nil {
		return nil, err
	}
	return nil, u.roleDal.Save(&role)
}

func (u *UserBll) deleteRole(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.roleDal.Delete(uId)
	return nil, err
}

func (u *UserBll) getRoleList(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}
