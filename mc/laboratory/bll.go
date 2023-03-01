/**
 * @Author: Joey
 * @Description:
 * @Create Date: 2023/2/21 16:38
 */

// Package laboratory
//  @Description: 检验信息
package laboratory

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

	labDal     *LabDal
	checkInDal *CheckInDal
	auditDal   *AuditDal
	resultDal  *ResultDal
	graphDal   *GraphDal
	reportDal  *ReportDal
	patientBll *patient.Bll
}

func (bll *Bll) RegApi(apiMap qf.ApiMap) {

	//apiMap.Reg(qf.EKindSave, "samples", bll.sampleDal.Save)
	//apiMap.Reg(qf.EKindDelete, "samples", bll.sampleDal.Delete)
	//
	//apiMap.Reg(qf.EKindSave, "labs", bll.labDal.Save)
	//apiMap.Reg(qf.EKindDelete, "labs", bll.labDal.Delete)
	//apiMap.Reg(qf.EKindSave, "labs/checkin", bll.checkin)
	//apiMap.Reg(qf.EKindDelete, "labs/audit", bll.audit)
	//
	//apiMap.Reg(qf.EKindSave, "results", bll.resultDal.Save)
	//apiMap.Reg(qf.EKindDelete, "results", bll.resultDal.Delete)
	//
	//apiMap.Reg(qf.EKindSave, "graphs", bll.graphDal.Save)
	//apiMap.Reg(qf.EKindDelete, "graphs", bll.graphDal.Delete)
	//
	//apiMap.Reg(qf.EKindSave, "reports", bll.reportDal.Save)
	//apiMap.Reg(qf.EKindDelete, "reports", bll.reportDal.Delete)
}

func (bll *Bll) RegDal(dalMap qf.DalMap) {
	//bll.orderDal =&OrderDal{}
	//dalMap.Reg(bll.orderDal, Order{})
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

//
// checkin
//  @Description: 上机
//  @receiver bll
//  @param ctx
//  @return interface{}
//  @return error
//
func (bll *Bll) checkin(ctx *qf.Context) (interface{}, error) {
	//上机时要求提交的是 样本条码号 日期 样本号
	//上机的本质 是save一条记录到labSample表
	//业务流程:
	//1 通过条码号先找到对应的样本记录
	//2 通过样本记录上的患者信息找到病例、患者信息并将病例、患者信息合并到检验样本信息表中保存
	//3 返回完整的检验样本信息供前端显示

	//	barcode :=ctx.GetValue("barcode")
	//	sampleNo :=ctx.GetValue("sample_no")
	//	SampleDate :=ctx.GetValue("sample_date")
	//	labId :=ctx.GetId("lab_id")
	//	sample, err := bll.sampleDal.GetFromBarcode(barcode)
	//	if err != nil {
	//		return nil, err
	//	}
	//	lab,err := bll.labDal.GetModel(labId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	if lab ==nil {
	//	b,e:=	bll.labDal.Save(lab)
	//	}
	//	labSample :=&LabSample{
	//	ID:lab.ID,
	//	content :ctx.ToContent(),
	//}
	//	ctx.BindModel(lab)
	//lab.
	//	model := &CheckIn{}
	//	model.PersonId = ctx.UserId
	//	model.LabId = ctx.UserId             //来自上下文
	//	model.Content = bll.ToContent(model) //转换为json
	//	rs, err := bll.checkInDal.Save(model)
	//	if rs == false {
	//		return nil, err
	//	}
	//	return model, nil
	return nil, nil
}

func (bll *Bll) audit(ctx *qf.Context) (interface{}, error) {
	return nil, nil
	//return bll.auditDal.Save(content)
}
