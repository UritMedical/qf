/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:52
 */

package liteLab

import (
	"qf"
)

//
// LabResult
//  @Description: 检验结果
//
type LabResult struct {
	//
	// PId
	//  @Description: 患者唯一号
	//
	PId uint `gorm:"index" json:"p_id"`

	//
	// CHId CaseHistoryId
	//  @Description: 病历号
	//
	CHId uint `gorm:"index" json:"ch_id"`

	//
	// LabDate
	//  @Description: 检验日期
	//
	LabDate uint `gorm:"index" json:"lab_date"`

	//
	// ItemId
	//  @Description: 检验项目
	//
	ItemId uint `gorm:"index" json:"item_id"`

	//
	// Value
	//  @Description: 检测值
	//
	Value uint `json:"value"`

	qf.Content
}
