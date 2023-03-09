/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/27 8:13
 */

package order

import "github.com/Urit-Mediacal/qf"

type ODal struct {
	qf.BaseDal
}

//
// SampleDal
//  @Description: 样本dal
//
type SampleDal struct {
	qf.BaseDal
}

func (dal *SampleDal) SaveBatch(batch []Sample) error {
	db := dal.DB()
	db.Save(batch)
	return db.Error
}
