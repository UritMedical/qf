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

func (dal *InfoDal) GetListByKey(key string, dest interface{}) qf.IError {
	err := dal.DB().Where("HisId = ? or Name LIKE ?", key, "%"+key+"%").Find(dest).Error
	if err != nil {
		return qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return nil
}

//
// DeleteByPatientId
//  @Description: 根据患者唯一号删除所有病历
//  @param pid
//  @return error
//
func (dal *CaseDal) DeleteByPatientId(pid uint64) qf.IError {
	rs := dal.DB().Where("PId = ?", pid).Delete(PatientCase{PId: pid})
	if rs.Error != nil {
		return qf.Error(qf.ErrorCodeDeleteFailure, rs.Error.Error())
	}
	return nil
}

func (dal *CaseDal) GetListByPatientId(pid uint64, dest interface{}) qf.IError {
	err := dal.DB().Where("PId = ?", pid).Find(dest).Error
	if err != nil {
		return qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return nil
}

func (dal *CaseDal) GetListByCaseId(caseId string, dest interface{}) qf.IError {
	err := dal.DB().Where("CaseId = ?", caseId).Find(dest).Error
	if err != nil {
		return qf.Error(qf.ErrorCodeRecordNotFound, err.Error())
	}
	return nil
}
