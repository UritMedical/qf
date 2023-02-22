package patient

import "qf"

type Patient struct {
	qf.Content
	HisId   string `gorm:"index"` // 患者的院内唯一号（来至于第三方系统，如HIS、体检等）
	PatNo   string `gorm:"index"` // 患者编号（住院号/门诊号/体检号等）
	PatName string `gorm:"index"` // 患者姓名
}
