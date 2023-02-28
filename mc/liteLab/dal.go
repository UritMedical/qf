/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/26 9:22
 */

package liteLab

import "qf"

//
// LabResultDal
//  @Description: 检验结果dal
//
type LabResultDal struct {
	qf.BaseDal
}

func (dal *LabResultDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (dal *LabResultDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}
