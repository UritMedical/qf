/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:52
 */

package liteLab

import (
	"github.com/Urit-Mediacal/qf"
)

//EValueFlag 结果标志枚举
type EValueFlag rune

const (
	//EValueFlagHigh  高值标志
	EValueFlagHigh = 'H'
	//EValueFlagLow  低值标志
	EValueFlagLow = 'L'
	//EValueFlagNormal  正常值标志
	EValueFlagNormal = 'N'
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
	PId uint `gorm:"index"`

	//
	// CId CaseId
	//  @Description: 病历号
	//
	CId uint `gorm:"index"`

	//
	// LabDate
	//  @Description: 检验日期
	//
	LabDate uint `gorm:"index"`

	//
	// ItemId
	//  @Description: 检验项目
	//
	ItemId uint `gorm:"index"`

	//
	// Value
	//  @Description: 检测值
	//
	Value uint

	//
	// Flag
	//  @Description: 结果标志
	//
	Flag EValueFlag

	qf.BaseModel
}
