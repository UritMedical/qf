package user

import "qf"

func (u *UserBll) regRightsApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "rights", u.saveRightGroup)
	api.Reg(qf.EKindDelete, "rights", u.deleteRightGroup)
	api.Reg(qf.EKindGetModel, "rights", u.getRightGroupList)
}

func (u *UserBll) saveRightGroup(ctx *qf.Context) (interface{}, error) {
	var rg RightsGroup
	if err := ctx.BindModel(&rg); err != nil {
		return nil, err
	}
	return nil, u.roleDal.Save(&rg)
}

func (u *UserBll) deleteRightGroup(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.roleDal.Delete(uId)
	return nil, err
}

func (u *UserBll) getRightGroupList(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}
