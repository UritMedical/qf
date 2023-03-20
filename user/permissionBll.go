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

func (b *Bll) savePermission(ctx *qf.Context) (interface{}, qf.IError) {
	permission := &model.Permission{}
	if err := ctx.Bind(permission); err != nil {
		return nil, err
	}
	if permission.Id == 0 {
		permission.Id = ctx.NewId(permission)
	}
	return nil, b.permissionDal.Save(permission)
}

func (b *Bll) deletePermission(ctx *qf.Context) (interface{}, qf.IError) {
	uId := ctx.GetId()
	return nil, b.permissionDal.Delete(uId)
}

func (b *Bll) getPermissions(ctx *qf.Context) (interface{}, qf.IError) {
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
func (b *Bll) setPermissionApi(ctx *qf.Context) (interface{}, qf.IError) {
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
func (b *Bll) getPermissionApi(ctx *qf.Context) (interface{}, qf.IError) {
	permissionId := ctx.GetUIntValue("PermissionId")
	return b.permissionApiDal.GetApisByPermissionId(permissionId)
}
