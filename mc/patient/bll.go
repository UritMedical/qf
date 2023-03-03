package patient

import (
	"qf"
)

type Bll struct {
	qf.BaseBll
	infoDal *InfoDal
	caseDal *CaseDal

	getDict qf.ApiHandler
}

func (b *Bll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "", b.SavePatient)      // 保存患者基本信息
	api.Reg(qf.EKindDelete, "", b.DeletePatient)  // 删除患者，包含基本信息和全部病历
	api.Reg(qf.EKindSave, "case", b.SaveCase)     // 保存患者病历信息
	api.Reg(qf.EKindDelete, "case", b.DeleteCase) // 删除单个病历
	api.Reg(qf.EKindGetModel, "", b.GetFull)      // 按唯一号或HIS唯一号获取完整信息
	api.Reg(qf.EKindGetList, "", b.GetFullList)   // 按条件获取完整列表
}

func (b *Bll) RegDal(dal qf.DalMap) {
	b.infoDal = &InfoDal{}
	b.caseDal = &CaseDal{}
	dal.Reg(b.infoDal, Patient{})
	dal.Reg(b.caseDal, Case{})
}

func (b *Bll) RegMsg(msg qf.MessageMap) {

}

func (b *Bll) RegRef(ref qf.RefMap) {

}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}

//
// SavePatient
//  @Description: 包含患者基本信息
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
	id := ctx.GetUIntValue("id")
	// 删除患者信息
	err := b.infoDal.Delete(id)
	if err == nil {
		// 删除所有病历
		err = b.caseDal.DeleteByPatientId(id)
	}
	return nil, err
}

func (b *Bll) SaveCase(ctx *qf.Context) (interface{}, error) {
	model := &Case{}
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

func (b *Bll) DeleteCase(ctx *qf.Context) (interface{}, error) {
	return nil, b.caseDal.Delete(ctx.GetUIntValue("id"))
}

func (b *Bll) GetFull(ctx *qf.Context) (interface{}, error) {
	id := ctx.GetUIntValue("id")

	// 通过ID检索
	model := &Patient{}
	err := b.infoDal.GetModel(id, model)
	if err != nil {
		return nil, err
	}
	return b.Map(model), nil
}

func (b *Bll) GetFullList(ctx *qf.Context) (interface{}, error) {

	return nil, nil
}

////
//// SavePatient
////  @Description: 保存用户基本信息
////  @param ctx
////  @return interface{}
////  @return error
////
//func (b *Bll) SavePatient(ctx *qf.Context) (interface{}, error) {
//	// 获取提交的基本信息
//	model := &Patient{}
//	err := ctx.Bind(model)
//	if err != nil {
//		return nil, err
//	}
//	// 刷新内容
//	model.Content = b.BuildContent(model)
//	// 保存患者基本信息
//	return model, b.infoDal.Save(model)
//}
//
//func (b *Bll) SaveCase(ctx *qf.Context) (interface{}, error) {
//	return nil, nil
//}
//
//func (b *Bll) DeleteCase(ctx *qf.Context) (interface{}, error) {
//	return nil, nil
//}
//
//func (b *Bll) DeletePatient(ctx *qf.Context) (interface{}, error) {
//	id := ctx.GetUIntValue("id")
//	return id, b.infoDal.Delete(id)
//}
//
////
//// GetCase
////  @Description: 通过病历唯一ID获取单条病历信息
////  @param ctx 传入结构体为：
////  @return interface{} 病历Case结构Json
////  @return error
////
//func (b *Bll) GetCase(ctx *qf.Context) (interface{}, error) {
//	// 外部调用时，请用框架创建上下文
//	// qf.BuildContext(map)
//
//	//id := ctx.GetUIntValue("id")
//	//b.caseDal.GetModel(id, Case{})
//	return nil, nil
//}
//
////
//// GetCaseList
////  @Description: 根据病历号或者患者唯一号，获取一组病历列表
////  @param ctx
////  @return interface{}
////  @return error
////
//func (b *Bll) GetCaseList(ctx *qf.Context) (interface{}, error) {
//	// 外部调用时，请用框架创建上下文
//	// qf.BuildContext(map)
//	return nil, nil
//	//infoId := ctx.GetUIntValue("infoId")
//	//caseId := ctx.GetStringValue("caseId")
//	//return b.caseDal.Search(infoId, caseId)
//}
//
////
//// GetPatientFull
////  @Description: 根据唯一号获取患者完整信息（基本+病历列表）
////  @param ctx
////  @return interface{}
////  @return error
////
//func (b *Bll) GetPatientFull(ctx *qf.Context) (interface{}, error) {
//	pkg := struct {
//		Patient         // 患者基本信息
//		CaseList []Case // 包含的病历列表
//	}{}
//	// 按患者唯一号，获取患者基本信息
//	id := ctx.GetUIntValue("id")
//	err := b.infoDal.GetModel(id, &pkg)
//	if pkg.Id == 0 {
//		return nil, nil
//	}
//	// 根据患者唯一号，获取所有病历列表
//	return pkg, err
//}
