/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/27 8:13
 */

package order

import "github.com/UritMedical/qf"

//
// Order
//  @Description: 医嘱 检验申请
//
type Order struct {

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

	qf.BaseModel
}

//
// ODetail
//  @Description: 医嘱明细
//
type ODetail struct {
	Id    uint
	Title string
}

//
// Sample
//  @Description: 标本 指的是实体标本的信息
//
type Sample struct {
	//
	// OrderId
	//  @Description: 医嘱号
	//
	OrderId uint `gorm:"index"`
	//
	// Barcode
	//  @Description: 条码号
	//
	Barcode string `gorm:"index"`

	//
	// content.Content
	//  @Description: ID 明文 分管后所包含的套餐明细、检验项目明细等
	//
	qf.BaseModel
}
