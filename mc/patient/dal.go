package patient

import "qf"

type InfoDal struct {
	qf.BaseDal
}

type CaseDal struct {
	qf.BaseDal
}

//
// DeleteByPatientId
//  @Description: 根据患者唯一号删除所有病历
//  @param pid
//  @return error
//
func (dal *CaseDal) DeleteByPatientId(pid uint64) error {
	return dal.DB().Where("PId = ?", pid).Delete(Case{PId: pid}).Error
}
