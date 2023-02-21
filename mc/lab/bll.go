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
	patientBll patient.Bll
}

func (bll *Bll) regApi(apiMap qf.ApiMap) {
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

func (bll *Bll) regDal(dalMap qf.DalMap) {
	dalMap.Reg(bll.orderDal, Order{})
}

func (bll *Bll) refBll() []qf.IBll {
	return []qf.IBll{bll.patientBll}
}

func (bll *Bll) init() error {
	return nil
}

func (bll *Bll) stop() {
	return
}

func (bll *Bll) checkin(context Context) (interface{}, error) {
	model :=&CheckIn{} //
	bll.GetModelFromJson(context,model)
	model.PersonId =context.UserId,//来自上下文
	model.Content = bll.GetContent(model)//转换为json
	return bll.checkInDal.Save(model)
}

func (bll *Bll) audit(content interface{}) (interface{}, error) {
	return bll.auditDal.Save(content)
}
