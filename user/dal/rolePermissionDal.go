package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
)

type RolePermissionDal struct {
	qf.BaseDal
}

//
// SetRolePermission
//  @Description: 给指定角色分配权限。先删除roleId所有的权限，然后再重新添加
//  @param roleId
//  @param permissionIds
//  @return error
//
func (r RolePermissionDal) SetRolePermission(roleId uint64, permissionIds []uint64) qf.IError {
	tx := r.DB().Begin()
	//先删除原来的权限
	if err := tx.Where("RoleId = ?", roleId).Delete(&model.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
	}

	//再将权限添加到数据库
	list := make([]model.RolePermission, 0)
	for _, id := range permissionIds {
		list = append(list, model.RolePermission{
			RoleId:       roleId,
			PermissionId: id,
		})
	}

	if err := tx.Create(&list).Error; err != nil {
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
// GetRolePermission
//  @Description: 获取指定角色拥有的权限
//  @param roleId
//
func (r RolePermissionDal) GetRolePermission(roleId uint64) ([]uint64, qf.IError) {
	permissions := make([]uint64, 0)
	err := r.DB().Where("RoleId = ?", roleId).Select("PermissionId").Find(&permissions).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return permissions, nil
}
