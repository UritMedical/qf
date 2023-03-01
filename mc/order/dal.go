/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/27 8:13
 */

package order

import "qf"

type OrderDal struct {
	qf.BaseDal
}

func (dal *OrderDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (dal *OrderDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

//
// SampleDal
//  @Description: 样本dal
//
type SampleDal struct {
	qf.BaseDal
}

func (dal *SampleDal) BeforeAction(kind qf.EKind, content interface{}) error {
	return nil
}

func (dal *SampleDal) AfterAction(kind qf.EKind, content interface{}) error {
	return nil
}

type LabDal struct {
	qf.BaseDal
}

func (dal *SampleDal) SaveBatch(batch []Sample) error {
	db := dal.DB()
	db.Save(batch)
	return db.Error
}
