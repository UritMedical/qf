/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/26 9:28
 */

package liteLab

import "github.com/Urit-Mediacal/qf"

type Bll struct {
	qf.BaseBll
	resultDal *LabResultDal
}

func (bll *Bll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EApiKindSave, "results", bll.labResultSave)
	api.Reg(qf.EApiKindDelete, "results", bll.labResultDelete)

}

func (bll *Bll) RegDal(dal qf.DalMap) {
	dal.Reg(bll.resultDal, LabResult{})
}

func (bll *Bll) RegMsg(msg qf.MessageMap) {
	//TODO implement me
}

func (bll *Bll) RegRef(ref qf.RefMap) {
	//TODO implement me
	//ref.Reg("patient", qf.EApiKindGetModel, "", bll.patientGetModel)
	//ref.Reg("labItem", qf.EApiKindGetModel, "", bll.labItemGetModel)
}

func (bll *Bll) Init() error {
	return nil
}

func (bll *Bll) Stop() {
}

//
// labResultSave
//  @Description: 提交结果 注意,前端务必提供 LabResult对象中的全部索引字段和用于显示的相关明文
//  @receiver bll
//  @param ctx
//  @return interface{}
//  @return error
//
func (bll *Bll) labResultSave(ctx *qf.Context) (interface{}, error) {

	model := &LabResult{}
	err := ctx.Bind(model)
	if err != nil {
		return nil, err
	}
	err = bll.resultDal.Save(model)
	return model, err
}

func (bll *Bll) labResultDelete(ctx *qf.Context) (interface{}, error) {
	return bll.resultDal.Delete(ctx.GetId())
}
