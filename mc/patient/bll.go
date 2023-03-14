package patient

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/util"
)

type Bll struct {
	qf.BaseBll
	infoDal *InfoDal
	caseDal *CaseDal
}

func (b *Bll) RegApi(a qf.ApiMap) {
	a.Reg(qf.EApiKindSave, "patient", b.SavePatient)       // 保存患者基本信息
	a.Reg(qf.EApiKindDelete, "patient", b.DeletePatient)   // 删除患者，包含基本信息和全部病历
	a.Reg(qf.EApiKindSave, "patient/case", b.SaveCase)     // 保存患者病历信息
	a.Reg(qf.EApiKindDelete, "patient/case", b.DeleteCase) // 删除单个病历
	a.Reg(qf.EApiKindGetModel, "patient", b.GetFull)       // 按唯一号或HIS唯一号获取完整信息（基本信息+病历列表）
	a.Reg(qf.EApiKindGetList, "patients", b.GetFullList)   // 按条件获取完整列表
}

func (b *Bll) RegDal(d qf.DalMap) {
	b.infoDal = &InfoDal{}
	b.caseDal = &CaseDal{}
	d.Reg(b.infoDal, Patient{})
	d.Reg(b.caseDal, PatientCase{})
}

func (b *Bll) RegMsg(_ qf.MessageMap) {

}

func (b *Bll) RegRef(_ qf.RefMap) {

}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}

//
// SavePatient
//  @Description: 保存患者基本信息
//  @param ctx 输入结构
//		{
//			"Id": 0,
//			"HisId": "院内唯一号",
//			"Name": "患者姓名",
//			"Sex": "患者性别",
//			"Birth": "出生日期",
//			"Phone": "联系电话",
//			"IDCard": "身份证号"
//		}
//  @return interface{} 患者唯一号
//  @return error 异常
func (b *Bll) SavePatient(ctx *qf.Context) (interface{}, error) {
	model := &Patient{}
	if err := ctx.Bind(model); err != nil {
		return nil, err
	}
	// 将空字符串作为nil
	if *model.HisId == "" {
		model.HisId = nil
	}

	// 提交，如果HisId重复，则返回失败
	err := b.infoDal.Save(model)
	if err != nil {
		return 0, err
	}
	return model.Id, err
}

//
// DeletePatient
//  @Description: 删除患者信息
//  @param ctx
//		{
//			"Id": 需要删除的患者唯一Id,
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) DeletePatient(ctx *qf.Context) (interface{}, error) {
	// 删除患者信息
	err := b.infoDal.Delete(ctx.GetId())
	if err == nil {
		// 删除所有病历
		err = b.caseDal.DeleteByPatientId(ctx.GetId())
	}
	return nil, err
}

//
// SaveCase
//  @Description: 保存病历
//  @param ctx
//		{
//			"Id": 0,
//			"PId": 患者唯一号,
//			"CaseId": "病历号（门诊号/住院号/体检号）",
//			"Classify": "分类（门诊/住院/体检）",
//			"MedHistory": "病史"
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) SaveCase(ctx *qf.Context) (interface{}, error) {
	model := &PatientCase{}
	if err := ctx.Bind(model); err != nil {
		return nil, err
	}
	// 提交，如果CaseId重复，则返回失败
	err := b.caseDal.Save(model)
	if err != nil {
		return 0, err
	}
	return model.Id, err
}

//
// DeleteCase
//  @Description: 删除一条病历
//  @param ctx
//		{
//			"Id": 需要删除的病历唯一Id,
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) DeleteCase(ctx *qf.Context) (interface{}, error) {
	return nil, b.caseDal.Delete(ctx.GetId())
}

//
// GetFull
//  @Description: 获取单个患者完整信息（基本信息+病历列表）
//  @param ctx
//		{
//			"Id": 患者唯一Id,
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) GetFull(ctx *qf.Context) (interface{}, error) {
	// 通过ID检索
	patInfo := Patient{}
	err := b.infoDal.GetModel(ctx.GetId(), &patInfo)
	if err != nil || patInfo.Id == 0 {
		return nil, err
	}
	// 通过患者Id获取所有病历
	caseList := make([]PatientCase, 0)
	err = b.caseDal.GetListByPatientId(patInfo.Id, &caseList)
	if err != nil {
		return nil, err
	}

	// 返回
	rt := struct {
		Patient interface{}
		Cases   interface{}
	}{
		Patient: util.ToMap(patInfo),
		Cases:   util.ToMaps(caseList),
	}
	return rt, nil
}

//
// GetFullList
//  @Description: 按条件获取完整列表 ?key=xxx
//  @param ctx
//		{
//			"key": "查询关键字，姓名、病历号、院内HIS唯一号",
//		}
//  @return interface{}
//  @return error
//
func (b *Bll) GetFullList(ctx *qf.Context) (interface{}, error) {
	key := ctx.GetStringValue("key")
	// 先查患者基本信息列表
	pats := make([]Patient, 0)
	err := b.infoDal.GetListByKey(key, &pats)
	if err != nil {
		return nil, err
	}

	var rts []struct {
		Patient interface{}
		Cases   interface{}
	}
	if len(pats) == 0 {
		// 基本信息未查询到数据，尝试查询病历表
		caseList := make([]PatientCase, 0)
		err = b.caseDal.GetListByCaseId(key, &caseList)
		if err != nil {
			return nil, err
		}
		// 遍历查询
		for _, c := range caseList {
			p := Patient{}
			err = b.infoDal.GetModel(c.PId, &p)
			if err != nil {
				return nil, err
			}
			rts = append(rts, struct {
				Patient interface{}
				Cases   interface{}
			}{
				Patient: util.ToMap(p),
				Cases:   util.ToMaps(caseList),
			})
		}
	} else {
		// 遍历查询
		for _, p := range pats {
			caseList := make([]PatientCase, 0)
			err = b.caseDal.GetListByPatientId(p.Id, &caseList)
			if err != nil {
				return nil, err
			}
			rts = append(rts, struct {
				Patient interface{}
				Cases   interface{}
			}{
				Patient: util.ToMap(p),
				Cases:   util.ToMaps(caseList),
			})
		}
	}
	return rts, nil
}
