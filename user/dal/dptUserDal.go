package dal

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/user/model"
	"github.com/UritMedical/qf/util"
)

type DptUserDal struct {
	qf.BaseDal
}

//
// AddUsers
//  @Description: 向指定部门添加用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) AddUsers(departId uint64, userIds []uint64) qf.IError {
	oldUserIds, err := d.GetUsersByDptId(departId)
	if err != nil {
		return err
	}

	//过滤出部门中已经存在的账号
	newUsers := util.DiffIntSet(userIds, oldUserIds)

	list := make([]model.DepartUser, 0)
	for _, id := range newUsers {
		list = append(list, model.DepartUser{
			DepartId: departId,
			UserId:   id,
		})
	}
	e := d.DB().Create(list).Error
	if e != nil {
		return qf.Error(qf.ErrorCodeSaveFailure, e.Error())
	}
	return nil
}

//
// RemoveUser
//  @Description: 从指定部门移除用户
//  @param departId
//  @param userIds
//  @return error
//
func (d DptUserDal) RemoveUser(departId uint64, userId uint64) qf.IError {
	err := d.DB().Where("DepartId = ? AND UserId = ?", departId, userId).Delete(&model.DepartUser{}).Error
	if err != nil {
		return qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
	}
	return nil
}

//
// GetUsersByDptId
//  @Description: 获取部门中所有用户
//  @param departId
//  @return []uint64
//  @return error
//
func (d DptUserDal) GetUsersByDptId(departId uint64) ([]uint64, qf.IError) {
	userIds := make([]uint64, 0)
	err := d.DB().Where("DepartId = ?", departId).Select("UserId").Find(&userIds).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return userIds, nil
}

//
// GetDptsByUserId
//  @Description: 获取用户所属部门
//  @param userId
//  @return []uint64
//  @return error
//
func (d DptUserDal) GetDptsByUserId(userId uint64) ([]uint64, qf.IError) {
	dptIds := make([]uint64, 0)
	err := d.DB().Where("UserId = ?", userId).Select("DepartId").Find(&dptIds).Error
	if err != nil {
		return nil, qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return dptIds, nil
}
