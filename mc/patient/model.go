package patient

import "qf"

//
// Info
//  @Description: 患者基础信息
//
type Info struct {
	//
	// PatName
	//  @Description: 患者姓名，索引
	//
	PatName string `gorm:"index"`

	//
	// HisId
	//  @Description: HIS唯一号，唯一索引，通过院内唯一号快速查找患者信息
	//
	HisId string `gorm:"uniqueIndex"`

	//
	// qf.Content
	//  @Description: 完整信息
	//
	qf.Content
}

//
// Case
//  @Description: 患者病历信息
//
type Case struct {
	//
	// BaseId
	//  @Description: 对应基本信息ID
	//
	InfoId uint `gorm:"index"`

	//
	// CaseId
	//  @Description: 病历号（门诊号/住院号）
	//
	CaseId string `gorm:"uniqueIndex"`

	//
	// Classify
	//  @Description: 分类（门诊病历/住院病历）
	//
	Classify string

	//
	// qf.Content
	//  @Description: 病历完整内容
	//
	qf.Content
}
