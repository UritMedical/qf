package user

import (
	"qf"
	uModel "qf/mc/user/model"
)

func (b *Bll) regRightsApi(api qf.ApiMap) {
	//权限组
	api.Reg(qf.EApiKindSave, "rights", b.saveRightsGroup)       //添加权限组
	api.Reg(qf.EApiKindDelete, "rights", b.deleteRightsGroup)   //删除权限组
	api.Reg(qf.EApiKindGetList, "rights", b.getRightsGroupList) //获取权限组

	//权限组-API
	api.Reg(qf.EApiKindSave, "rights/apis", b.setRightsGroupApi)
	api.Reg(qf.EApiKindGetList, "rights/apis", b.getRightsGroupApi)
}

func (b *Bll) saveRightsGroup(ctx *qf.Context) (interface{}, error) {
	var rg uModel.RightsGroup
	if err := ctx.Bind(&rg); err != nil {
		return nil, err
	}
	return nil, b.rightsDal.Save(&rg)
}

func (b *Bll) deleteRightsGroup(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetId()
	ret, err := b.rightsDal.Delete(uId)
	return ret, err
}

func (b *Bll) getRightsGroupList(ctx *qf.Context) (interface{}, error) {
	rights := make([]uModel.RightsGroup, 0)
	err := b.rightsDal.GetList(0, 100, &rights)
	return b.Maps(rights), err
}

//
// setRightsGroupApi
//  @Description: 批量设置权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) setRightsGroupApi(ctx *qf.Context) (interface{}, error) {
	params := struct {
		RightsId uint64
		ApiIds   []string
	}{}
	if err := ctx.Bind(&params); err != nil {
		return nil, err
	}
	return nil, b.rightsApiDal.SetRightsApis(params.RightsId, params.ApiIds)
}

//
// getRightsGroupApi
//  @Description: 获取指定权限组能访问的API
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) getRightsGroupApi(ctx *qf.Context) (interface{}, error) {
	rightId := ctx.GetUIntValue("RightsId")
	return b.rightsApiDal.GetApisByRightsId(rightId)
}
