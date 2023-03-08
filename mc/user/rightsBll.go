package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (u *Bll) regRightsApi(api qf.ApiMap) {
	//权限组
	api.Reg(qf.EApiKindSave, "rights", u.saveRightsGroup)       //添加权限组
	api.Reg(qf.EApiKindDelete, "rights", u.deleteRightsGroup)   //删除权限组
	api.Reg(qf.EApiKindGetList, "rights", u.getRightsGroupList) //获取权限组

	//权限组-API
	api.Reg(qf.EApiKindSave, "rights/apis", u.setRightsGroupApi)
	api.Reg(qf.EApiKindGetList, "rights/apis", u.getRightsGroupApi)
}

func (u *Bll) saveRightsGroup(ctx *qf.Context) (interface{}, error) {
	var rg uModel.RightsGroup
	if err := ctx.Bind(&rg); err != nil {
		return nil, err
	}
	return nil, u.rightsDal.Save(&rg)
}

func (u *Bll) deleteRightsGroup(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetUIntValue("Id")
	ret, err := u.rightsDal.Delete(uId)
	return ret, err
}

func (u *Bll) getRightsGroupList(ctx *qf.Context) (interface{}, error) {
	rights := make([]uModel.RightsGroup, 0)
	err := u.rightsDal.GetList(0, 100, &rights)
	return u.Maps(rights), err
}

//
// setRightsGroupApi
//  @Description: 批量设置权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *Bll) setRightsGroupApi(ctx *qf.Context) (interface{}, error) {
	params := struct {
		RightsId uint64
		ApiIds   []string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, u.rightsApiDal.SetRightsApis(params.RightsId, params.ApiIds)
}

//
// getRightsGroupApi
//  @Description: 获取指定权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (u *Bll) getRightsGroupApi(ctx *qf.Context) (interface{}, error) {
	rightId := ctx.GetUIntValue("RightsId")
	return u.rightsApiDal.GetApisByRightsId(rightId)
}
