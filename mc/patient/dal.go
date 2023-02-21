package patient

import "qf"

type Dal struct {
	qf.BaseDal
}

func (d *Dal) BeforeAction(kind qf.EKind, content interface{}) (bool, error) {
	m := content.(Patient)
	// 检测重复
	if kind == qf.EKindSave {
		// 检测患者代号是否存在，存在返回false
		if d.DB().Where("pat_no = ?", m.PatNo).RowsAffected > 0 {
			return false, nil
		}
	}
	return true, nil
}

func (d *Dal) AfterAction(kind qf.EKind, content interface{}) (bool, error) {
	return true, nil
}

func (d *Dal) Search(content interface{}) (interface{}, error) {
	//d.DB().Where(...)
	return nil, nil
}
