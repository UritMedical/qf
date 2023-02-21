/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:52
 */

package lab

import (
	"qf/helper/content"
)

//
// Order
//  @Description: 医嘱 检验申请
//
type Order struct {
	//
	// CHId CaseHistoryId
	//  @Description: 病历号
	//
	CHId uint
	//
	// DocId
	//  @Description: 送检医生
	//
	DocId uint

	//
	// content.Content
	//  @Description:ID 摘要 明文套餐明细等
	//
	content.Content
}

//
// Sample
//  @Description: 标本
//
type Sample struct {
	//
	// OrderId
	//  @Description: 医嘱号
	//
	OrderId uint
	//
	// Barcode
	//  @Description: 条码号
	//
	Barcode string

	//
	// content.Content
	//  @Description: ID 明文 分管后所包含的套餐明细、检验项目明细等
	//
	content.Content
}

//
// Laboratory
//  @Description: 检验
//
type Laboratory struct {

	//
	// SampleId
	//  @Description: 样本号
	//
	SampleId uint
	//
	// SampleNo
	//  @Description: 标本号
	//
	SampleNo string `gorm:"index" json:"sample_no"`

	content.Content
}

//
// CheckIn
//  @Description: 上机
//
type CheckIn struct {

	//
	// LaboratoryId
	//  @Description: 检验活动Id
	//
	LaboratoryId uint
	//
	// PersonId
	//  @Description: 上机人id
	//
	PersonId uint
	//
	// content.Content
	//  @Description: 明文时间等
	//
	content.Content
}

//
// Audit
//  @Description: 审核
//
type Audit struct {

	//
	// LaboratoryId
	//  @Description: 检验活动Id
	//
	LaboratoryId uint
	//
	// AuditorId
	//  @Description: 审核者id
	//
	AuditorId uint
	//
	// content.Content
	//  @Description: 明文时间等
	//
	content.Content
}

//
// Result
//  @Description: 检验结果
//
type Result struct {
	//
	// LaboratoryId
	//  @Description: 检验id
	//
	LaboratoryId uint `gorm:"index" json:"laboratory_id"`
	//
	// SampleId
	//  @Description: 样本id
	//
	SampleId uint `gorm:"index" json:"sample_id"`
	//
	// content.Content
	//  @Description: 应该是一个集合 一般来说应包含当次检验的所有项目，包括项目id 检验结果 参考范围 异常标志等
	//
	content.Content
}

//
// Graph
//  @Description: 检验图像
//
type Graph struct {
	//
	// LaboratoryId
	//  @Description: 检验id
	//
	LaboratoryId uint `gorm:"index" json:"laboratory_id"`
	//
	// SampleId
	//  @Description: 样本id
	//
	SampleId uint `gorm:"index" json:"sample_id"`
	//
	// content.Content
	//  @Description: 应该是一个集合 一般来说应包含当次检验的所有项目，包括项目id 检验结果 参考范围 异常标志等
	//
	content.Content
}

//
// Report
//  @Description: 检验报告
//
type Report struct {

	//
	// PatientId
	//  @Description: 患者号
	//
	PatientId uint `gorm:"index" json:"patient_id"`

	//
	// LaboratoryId
	//  @Description: 病历号
	//
	CHId uint `gorm:"index" json:"ch_id"`
	//
	// OrderId
	//  @Description: 申请号
	//
	OrderId uint `gorm:"index" json:"order_id"`

	//
	// content.Content
	//  @Description: 完整报告内容
	//
	content.Content
}
