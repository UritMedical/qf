/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:52
 */

package laboratory

import (
	"github.com/UritMedical/qf"
)

//
// Lab
//  @Description: 检验
//
type Lab qf.BaseModel

//
// LabSample
//  @Description:检验标本信息 上机时附加
//
type LabSample struct {
	//
	// SampleId
	//  @Description: 样本号
	//
	SampleId uint `gorm:"index" json:"sample_id"`
	//
	// SampleNo
	//  @Description: 标本号
	//
	SampleNo string `gorm:"index" json:"sample_no"`
	//
	// SampleDate
	//  @Description: 检验日期
	//
	SampleDate uint `gorm:"index" json:"sample_date"`
	qf.BaseModel
}

//
// LabAudit
//  @Description: 审核
//
type LabAudit struct {

	//
	// AuditorId
	//  @Description: 审核者id
	//
	UserId uint

	//
	// content.Content
	//  @Description: 明文时间等
	//
	qf.BaseModel
}

//
// LabResult
//  @Description: 检验结果
//
type LabResult struct {
	//
	// LaboratoryId
	//  @Description: 检验id
	//
	LabId uint `gorm:"index"`

	//
	// content.Content
	//  @Description: 应该是一个集合 一般来说应包含当次检验的所有项目，包括项目id 检验结果 参考范围 异常标志等
	//
	qf.BaseModel
}

//
// LabGraph
//  @Description: 检验图像
//
type LabGraph LabResult

//
// LabReport
//  @Description: 检验报告
//
type LabReport qf.BaseModel
