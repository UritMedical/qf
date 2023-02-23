package patient

import "qf"

type Patient struct {
	//
	//  PatName
	//  @Description: 患者姓名，索引
	//
	PatName string `gorm:"index"`

	//
	//  qf.Content
	//  @Description: 唯一号，其他信息
	//
	qf.Content
}
