package uDal

import (
	"qf"
	uModel "qf/mc/user/model"
)

type RoleRightsDal struct {
	qf.BaseDal
}

func (r RoleRightsDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (r RoleRightsDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// SetRoleRights
//  @Description: 给指定角色分配权限。先删除roleId所有的权限，然后再重新添加
//  @param roleId
//  @param rightsIds
//  @return error
//
func (r RoleRightsDal) SetRoleRights(roleId uint, rightsIds []uint) error {
	tx := r.DB().Begin()
	//先删除原来的权限
	if err := tx.Where("role_id = ?", roleId).Delete(&uModel.RoleRights{}).Error; err != nil {
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
func (r RoleRightsDal) GetRoleRights(roleId uint) ([]uint, error) {
	rights := make([]uint, 0)
	err := r.DB().Where("role_id = ?", roleId).Select("rights_id").Find(&rights).Error
	return rights, err
}
