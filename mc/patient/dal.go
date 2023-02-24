package patient

import "qf"

type InfoDal struct {
	qf.BaseDal
}

func (b *InfoDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (b *InfoDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (b *InfoDal) GetInfoByHisId(hisId string) (interface{}, error) {
	return nil, nil
}

type CaseBll struct {
	qf.BaseDal
}

func (c *CaseBll) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (c *CaseBll) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (c *CaseBll) Search(infoId uint, caseId string) (interface{}, error) {
	return nil, nil
}
