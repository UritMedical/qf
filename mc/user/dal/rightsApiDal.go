package uDal

import (
	"qf"
	uModel "qf/mc/user/model"
)

type RightsApiDal struct {
	qf.BaseDal
}

func (r RightsApiDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (r RightsApiDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// SetRightsApis
//  @Description: 向指定权限组添加，删除API
//  @param rightsId
//  @param apiKeys
//  @return error
//
func (r RightsApiDal) SetRightsApis(rightsId uint64, apiKeys []string) error {
	tx := r.DB().Begin()
	//先删除此权限组所有的API
	if err := tx.Where("RightsId = ?", rightsId).Delete(&uModel.RightsApi{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	apis := make([]uModel.RightsApi, 0)
	for _, key := range apiKeys {
		apis = append(apis, uModel.RightsApi{
			RightsId: rightsId,
			ApiId:    key,
		})
	}

	if err := tx.Create(&apis).Error; err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

//
// GetRightsApi
//  @Description: 获取指定权限组的API
//  @param rightsId
//  @return []string
//  @return error
//
func (r RightsApiDal) GetApisByRightsId(rightsId uint64) ([]string, error) {
	apis := make([]string, 0)
	err := r.DB().Where("RightsId = ?", rightsId).Select("ApiId").Find(&apis).Error
	return apis, err
}
