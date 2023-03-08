package uDal

import (
	"qf"
	"qf/mc/user/uModel"
)

type RoleRightsDal struct {
	qf.BaseDal
}

//
// SetRoleRights
//  @Description: 给指定角色分配权限。先删除roleId所有的权限，然后再重新添加
//  @param roleId
//  @param rightsIds
//  @return error
//
func (r RoleRightsDal) SetRoleRights(roleId uint64, rightsIds []uint64) error {
	tx := r.DB().Begin()
	//先删除原来的权限
	if err := tx.Where("RoleId = ?", roleId).Delete(&uModel.RoleRights{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	//再将权限添加到数据库
	list := make([]uModel.RoleRights, 0)
	for _, id := range rightsIds {
		list = append(list, uModel.RoleRights{
			RoleId:   roleId,
			RightsId: id,
		})
	}

	if err := tx.Create(&list).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

//
// GetRoleRights
//  @Description: 获取指定角色拥有的权限
//  @param roleId
//
func (r RoleRightsDal) GetRoleRights(roleId uint64) ([]uint64, error) {
	rights := make([]uint64, 0)
	err := r.DB().Where("RoleId = ?", roleId).Select("RightsId").Find(&rights).Error
	return rights, err
}
