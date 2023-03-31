/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 17:39
 */

package laboratory

import "github.com/UritMedical/qf"

//
// LabDal
//  @Description: 检验
//
type LabDal struct {
	qf.BaseDal
}

//
// CheckInDal
//  @Description: 上机dal
//
type CheckInDal struct {
	qf.BaseDal
}

//
// AuditDal
//  @Description: 审核dal
//
type AuditDal struct {
	qf.BaseDal
}

//
// ResultDal
//  @Description: 检验结果dal
//
type ResultDal struct {
	qf.BaseDal
}

//
// GraphDal
//  @Description: 检验图像dal
//
type GraphDal struct {
	qf.BaseDal
}

//
// ReportDal
//  @Description: 检验报告dal
//
type ReportDal struct {
	qf.BaseDal
}
