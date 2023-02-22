/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:38
 */

// Package laboratory
//  @Description: 检验信息
package lab

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
	labDal     *LaboratoryDal
	checkInDal *CheckInDal
	auditDal   *AuditDal
	resultDal  *ResultDal
	graphDal   *GraphDal
	reportDal  *ReportDal
	patientBll *patient.Bll
}

func (bll *Bll) RegApi(apiMap qf.ApiMap) {
	apiMap.Reg(qf.EKindSave, "orders", bll.orderDal.Save)
	apiMap.Reg(qf.EKindDelete, "orders", bll.orderDal.Delete)

	apiMap.Reg(qf.EKindSave, "samples", bll.sampleDal.Save)
	apiMap.Reg(qf.EKindDelete, "samples", bll.sampleDal.Delete)

	apiMap.Reg(qf.EKindSave, "labs", bll.labDal.Save)
	apiMap.Reg(qf.EKindDelete, "labs", bll.labDal.Delete)
	apiMap.Reg(qf.EKindSave, "labs/checkin", bll.checkin)
	apiMap.Reg(qf.EKindDelete, "labs/audit", bll.audit)

	apiMap.Reg(qf.EKindSave, "results", bll.resultDal.Save)
	apiMap.Reg(qf.EKindDelete, "results", bll.resultDal.Delete)

	apiMap.Reg(qf.EKindSave, "graphs", bll.graphDal.Save)
	apiMap.Reg(qf.EKindDelete, "graphs", bll.graphDal.Delete)

	apiMap.Reg(qf.EKindSave, "reports", bll.reportDal.Save)
	apiMap.Reg(qf.EKindDelete, "reports", bll.reportDal.Delete)
}

func (bll *Bll) RegDal(dalMap qf.DalMap) {
	dalMap.Reg(bll.orderDal, Order{})
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

func (bll *Bll) checkin(ctx *qf.Context) (interface{}, error) {
	model := &CheckIn{}
	ctx.BindModel(model)
	model.PersonId = ctx.UserId          //来自上下文
	model.Content = ctx.ToContent(model) //转换为json
	rs, err := bll.checkInDal.Save(model)
	if rs == false {
		return nil, err
	}
	return model, nil
}

func (bll *Bll) audit(ctx *qf.Context) (interface{}, error) {
	return nil, nil
	//return bll.auditDal.Save(content)
}
