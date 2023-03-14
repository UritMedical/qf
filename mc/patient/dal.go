package patient

import (
	"github.com/UritMedical/qf"
)

type InfoDal struct {
	qf.BaseDal
}

type CaseDal struct {
	qf.BaseDal
}

func (dal *InfoDal) GetListByKey(key string, dest interface{}) error {
	return dal.DB().Where("HisId = ? or Name LIKE ?", key, "%"+key+"%").Find(dest).Error
}

//
// DeleteByPatientId
//  @Description: 根据患者唯一号删除所有病历
//  @param pid
//  @return error
//
func (dal *CaseDal) DeleteByPatientId(pid uint64) (bool, error) {
	rs := dal.DB().Where("PId = ?", pid).Delete(PatientCase{PId: pid})
	return rs.RowsAffected > 0, rs.Error
}

func (dal *CaseDal) GetListByPatientId(pid uint64, dest interface{}) error {
	return dal.DB().Where("PId = ?", pid).Find(dest).Error
}

func (dal *CaseDal) GetListByCaseId(caseId string, dest interface{}) error {
	return dal.DB().Where("CaseId = ?", caseId).Find(dest).Error
}
