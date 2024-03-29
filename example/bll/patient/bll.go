package patient

import (
	"fmt"
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/util"
	"github.com/UritMedical/qf/util/qio"
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
	a.Reg(qf.EApiKindSave, "upload", b.uploadFile)         // 上传文件
}

func (b *Bll) RegDal(d qf.DalMap) {
	b.infoDal = &InfoDal{}
	b.caseDal = &CaseDal{}
	d.Reg(b.infoDal, Patient{})
	d.Reg(b.caseDal, nil)
}

func (b *Bll) RegFault(_ qf.FaultMap) {

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
func (b *Bll) SavePatient(ctx *qf.Context) (interface{}, qf.IError) {
	model := &Patient{}
	if err := ctx.Bind(model); err != nil {
		return nil, err
	}
	// 获取ID
	if model.Id == 0 {
		model.Id = ctx.NewId(model)
	}
	// 将空字符串作为nil
	if model.HisId != nil && *model.HisId == "" {
		model.HisId = nil
	}

	// 提交，如果HisId重复，则返回失败
	err := b.infoDal.Save(model)
	if err != nil {
		return 0, qf.Error(qf.ErrorCodeSaveFailure, err.Error())
	}
	return model.Id, nil
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
func (b *Bll) DeletePatient(ctx *qf.Context) (interface{}, qf.IError) {
	// 删除患者信息
	err := b.infoDal.Delete(ctx.GetId())
	if err == nil {
		// 删除所有病历
		err = b.caseDal.DeleteByPatientId(ctx.GetId())
	}
	return nil, qf.Error(qf.ErrorCodeDeleteFailure, err.Error())
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
func (b *Bll) SaveCase(ctx *qf.Context) (interface{}, qf.IError) {
	model := &PatientCase{}
	if err := ctx.Bind(model); err != nil {
		return nil, err
	}
	if model.Id == 0 {
		model.Id = ctx.NewId(model)
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
func (b *Bll) DeleteCase(ctx *qf.Context) (interface{}, qf.IError) {
	return nil, b.caseDal.Delete(ctx.GetId())
}

type AA struct {
	DT qf.Date
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
func (b *Bll) GetFull(ctx *qf.Context) (interface{}, qf.IError) {
	// 通过ID检索
	patInfo := Patient{}
	err := b.infoDal.GetModel(ctx.GetId(), &patInfo)
	if err != nil || patInfo.Id == 0 {
		return nil, err
	}

	// 返回
	rt := struct {
		Patient interface{}
	}{
		Patient: util.ToMap(patInfo),
	}
	return rt, nil
}

type Info struct {
	Patient
	Case  PatientCase
	Cases []PatientCase
	ABS   string
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
func (b *Bll) GetFullList(ctx *qf.Context) (interface{}, qf.IError) {
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

func (b *Bll) uploadFile(ctx *qf.Context) (interface{}, qf.IError) {
	files, e := ctx.GetFile("File")
	if e != nil {
		return nil, e
	}
	for _, file := range files {
		_ = qio.WriteAllBytes(fmt.Sprintf("D:/file/%s", file.Name), file.Data, false)
	}
	return nil, nil
}
