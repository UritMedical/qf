package patient

import "qf"

type BaseDal struct {
	qf.BaseDal
}

func (b *BaseDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (b *BaseDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}
