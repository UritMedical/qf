package patient

import "github.com/UritMedical/qf"

//
// Patient
//  @Description: 患者基础信息
//
type Patient struct {
	qf.BaseModel

	//
	// Name
	//  @Description: 患者姓名，索引
	//
	Name string `gorm:"index"`

	//
	// HisId
	//  @Description: HIS唯一号，唯一索引
	//
	HisId *string `gorm:"uniqueIndex:null"`
}

//
// PatientCase
//  @Description: 患者病历信息
//
type PatientCase struct {
	qf.BaseModel

	//
	// PId
	//  @Description: 对应基本信息ID
	//
	PId uint64 `gorm:"index"`

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
}
