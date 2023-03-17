package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
)

type PermissionApiDal struct {
	qf.BaseDal
}

//
// SetPermissionApis
//  @Description: 向指定权限组添加，删除API
//  @param permissionId
//  @param apiKeys
//  @return error
//
func (r PermissionApiDal) SetPermissionApis(permissionId uint64, apiKeys []string) qf.IError {
	tx := r.DB().Begin()
	//先删除此权限组所有的API
	if err := tx.Where("PermissionId = ?", permissionId).Delete(&model.PermissionApi{}).Error; err != nil {
		tx.Rollback()
		return qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
	}

	apis := make([]model.PermissionApi, 0)
	for _, key := range apiKeys {
		apis = append(apis, model.PermissionApi{
			PermissionId: permissionId,
			ApiId:        key,
		})
	}

	if err := tx.Create(&apis).Error; err != nil {
		tx.Rollback()
		return qf.Error(qf.ErrorCodeSaveFailure, err.Error())
	}
	e := tx.Commit().Error
	if e != nil {
		return qf.Error(qf.ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// GetApisByPermissionId
//  @Description: 获取指定权限组的API
//  @param permissionId
//  @return []string
//  @return error
//
func (r PermissionApiDal) GetApisByPermissionId(permissionId uint64) ([]string, qf.IError) {
	apis := make([]string, 0)
	err := r.DB().Where("PermissionId = ?", permissionId).Select("ApiId").Find(&apis).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return apis, nil
}
