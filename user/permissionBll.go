package user

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
	"github.com/UritMedical/qf/util"
)

func (b *Bll) regPermissionApi(api qf.ApiMap) {
	//权限组
	api.Reg(qf.EApiKindSave, "permission", b.savePermission)     //添加权限组
	api.Reg(qf.EApiKindDelete, "permission", b.deletePermission) //删除权限组
	api.Reg(qf.EApiKindGetList, "permissions", b.getPermissions) //获取权限组

	//权限组-API
	api.Reg(qf.EApiKindSave, "permission/apis", b.setPermissionApi)
	api.Reg(qf.EApiKindGetList, "permission/apis", b.getPermissionApi)
}

func (b *Bll) savePermission(ctx *qf.Context) (interface{}, error) {
	var rg model.Permission
	if err := ctx.Bind(&rg); err != nil {
		return nil, err
	}
	return nil, b.permissionDal.Save(&rg)
}

func (b *Bll) deletePermission(ctx *qf.Context) (interface{}, error) {
	uId := ctx.GetId()
	ret, err := b.permissionDal.Delete(uId)
	return ret, err
}

func (b *Bll) getPermissions(ctx *qf.Context) (interface{}, error) {
	permissions := make([]model.Permission, 0)
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
func (b *Bll) setPermissionApi(ctx *qf.Context) (interface{}, error) {
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
func (b *Bll) getPermissionApi(ctx *qf.Context) (interface{}, error) {
	permissionId := ctx.GetUIntValue("PermissionId")
	return b.permissionApiDal.GetApisByPermissionId(permissionId)
}
