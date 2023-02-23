package patient

import "qf"

//type Dal struct {
//	qf.BaseDal
//}
//
//func (d *Dal) BeforeAction(kind qf.EKind, content interface{}) error {
//	m := content.(Patient)
//	// 检测重复
//	if kind == qf.EKindSave {
//		// 检测患者代号是否存在，存在返回false
//		if d.DB().Where("pat_no = ?", m.PatNo).RowsAffected > 0 {
//			return errors.New("pat_no already exists")
//		}
//	}
//	return nil
//}
//
//func (d *Dal) AfterAction(kind qf.EKind, content interface{}) error {
//	return nil
//}

type Dal struct {
	qf.BaseDal
}

func (d *Dal) BeforeAction(kind qf.EKind, content interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (d *Dal) AfterAction(kind qf.EKind, content interface{}) error {
	//TODO implement me
	panic("implement me")
}
