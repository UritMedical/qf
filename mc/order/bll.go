/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/27 8:13
 */

package order

import (
	"github.com/UritMedical/qf"
)

//
// Bll
//  @Description: 业务逻辑
//
type Bll struct {
	qf.BaseBll

	orderDal  *ODal
	sampleDal *SampleDal
}

func (bll *Bll) RegMsg(_ qf.MessageMap) {

}

func (bll *Bll) RegRef(_ qf.RefMap) {
}

func (bll *Bll) Init() error {
	return nil
}

func (bll *Bll) Stop() {
}

func (bll *Bll) RegApi(apiMap qf.ApiMap) {
	apiMap.Reg(qf.EApiKindSave, "", bll.orderSave)
	apiMap.Reg(qf.EApiKindSave, "", bll.orderDel)

	apiMap.Reg(qf.EApiKindSave, "samples", bll.sampleSaveBatch)
	apiMap.Reg(qf.EApiKindSave, "samples", bll.sampleDel)

}

func (bll *Bll) RegDal(dalMap qf.DalMap) {
	dalMap.Reg(bll.orderDal, Order{})
	dalMap.Reg(bll.sampleDal, Sample{})
}

func (bll *Bll) orderSave(ctx *qf.Context) (interface{}, error) {
	model := Order{}
	err := ctx.Bind(model)
	if err != nil {
		return nil, err
	}
	return model, bll.orderDal.Save(model)
}

func (bll *Bll) orderDel(ctx *qf.Context) (interface{}, error) {
	return bll.orderDal.Delete(ctx.GetId())
}

func (bll *Bll) sampleSaveBatch(ctx *qf.Context) (interface{}, error) {
	var list []Sample
	err := ctx.Bind(list)
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
	return bll.sampleDal.Delete(ctx.GetId())
}
