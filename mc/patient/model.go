package patient

import "qf"

//
// Base
//  @Description: 患者基础信息
//
type Base struct {
	//
	// PatName
	//  @Description: 患者姓名，索引
	//
	PatName string `gorm:"index"`

	//
	// qf.Content
	//  @Description: 其他基本信息，如性别、生日、联系方式等
	//
	qf.Content
}

//
// Number
//  @Description: 患者院内号，此表的主要作用为通过院内号码快速定位的患者编号
//
type Number struct {
	//
	// HisId
	//  @Description: 院内唯一号
	//
	HisId string `gorm:"index"`

	//
	// CardId
	//  @Description: 院内就诊号（门诊号/住院号/体检号等）
	//
	CardId string `gorm:"index"`

	//
	// BaseId
	//  @Description: 对应基本信息ID
	//
	BaseId uint
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
	BaseId uint `gorm:"index"`

	//
	// CaseId
	//  @Description: 本次就诊流水号
	//
	CaseId string `gorm:"index"`

	//
	// From
	//  @Description: 本次就诊方式（门诊/住院/体检等）
	//
	From string

	//
	// qf.Content
	//  @Description: 其他基本信息，如就诊科室、主治医生、病史等
	//
	qf.Content
}
