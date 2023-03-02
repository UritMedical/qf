package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (u *UserBll) regRightsApi(api qf.ApiMap) {
	//权限组
	api.Reg(qf.EKindSave, "rights", u.saveRightGroup)        //添加权限组
	api.Reg(qf.EKindDelete, "rights", u.deleteRightGroup)    //删除权限组
	api.Reg(qf.EKindGetModel, "rights", u.getRightGroupList) //获取权限组

	//权限组-API
	api.Reg(qf.EKindSave, "rights/apis", u.setRightGroupApi)
}

func (u *UserBll) saveRightGroup(ctx *qf.Context) (interface{}, error) {
	var rg uModel.RightsGroup
	if err := ctx.Bind(&rg); err != nil {
		return nil, err
	}
	return nil, u.rightsDal.Save(&rg)
}

func (u *UserBll) deleteRightGroup(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	err := u.rightsDal.Delete(uId)
	return nil, err
}

func (u *UserBll) getRightGroupList(ctx *qf.Context) (interface{}, error) {
	rights := make([]uModel.RightsGroup, 0)
	err := u.rightsDal.GetList(0, 100, &rights)
	return rights, err
}

func (u *UserBll) setRightGroupApi(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}
