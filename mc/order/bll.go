/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/27 8:13
 */

package order

import (
	"qf"
	"qf/mc/patient"
)

//
// Bll
//  @Description: 业务逻辑
//
type Bll struct {
	qf.BaseBll

	orderDal   *OrderDal
	sampleDal  *SampleDal
	patientBll *patient.Bll
}

func (bll *Bll) RegApi(apiMap qf.ApiMap) {
	apiMap.Reg(qf.EKindSave, "orders", bll.orderSave)
	apiMap.Reg(qf.EKindDelete, "orders", bll.orderDel)

	apiMap.Reg(qf.EKindSave, "samples", bll.sampleSaveBatch)
	apiMap.Reg(qf.EKindDelete, "samples", bll.sampleDel)

}

func (bll *Bll) RegDal(dalMap qf.DalMap) {
	bll.orderDal = &OrderDal{}
	dalMap.Reg(bll.orderDal, Order{})
	dalMap.Reg(bll.sampleDal, Sample{})
}

func (bll *Bll) RefBll() []qf.IBll {
	return []qf.IBll{bll.patientBll}
}

func (bll *Bll) RegMsg(msg qf.MessageMap) {

}

func (bll *Bll) Init() error {
	return nil
}

func (bll *Bll) Stop() {
	return
}

func (bll *Bll) orderSave(ctx *qf.Context) (interface{}, error) {
	model := Order{}
	err := ctx.BindModel(model)
	if err != nil {
		return nil, err
	}
	return model, bll.orderDal.Save(model)
}

func (bll *Bll) orderDel(ctx *qf.Context) (interface{}, error) {
	id := ctx.GetUIntValue("id")
	return nil, bll.orderDal.Delete(id)
}

func (bll *Bll) sampleSaveBatch(ctx *qf.Context) (interface{}, error) {
	var list []Sample
	err := ctx.BindModel(list)
	if err != nil {
		return nil, err
	}
	return list, bll.sampleDal.SaveBatch(list)
}

//
// sampleDel
//  @Description: 删除样本
//  @receiver bll
//  @param ctx
//  @return interface{}
//  @return error
//
func (bll *Bll) sampleDel(ctx *qf.Context) (interface{}, error) {
	id := ctx.GetUIntValue("id")
	return nil, bll.sampleDal.Delete(id)
}
