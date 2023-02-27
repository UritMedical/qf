package patient

import (
	"qf"
)

type Bll struct {
	qf.BaseBll
	infoDal *InfoDal
	caseDal *CaseBll
}

func (b *Bll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "info", b.SaveInfo)     // 保存患者基本信息
	api.Reg(qf.EKindSave, "case", b.SaveCase)     // 保存患者病历
	api.Reg(qf.EKindDelete, "case", b.DeleteCase) // 删除患者病历
	api.Reg(qf.EKindDelete, "", b.DeletePatient)  // 删除患者

	api.Reg(qf.EKindGetModel, "case", b.GetCase)     // 通过病历唯一ID获取单条病历信息
	api.Reg(qf.EKindGetList, "cases", b.GetCaseList) // 根据病历号或者患者唯一号，获取一组病历列表
	api.Reg(qf.EKindGetModel, "", b.GetPatientInfo)  // 根据病历唯一ID或者HIS唯一号获取患者基本信息

}

func (b *Bll) RegDal(dal qf.DalMap) {
	b.infoDal = &InfoDal{}
	b.caseDal = &CaseBll{}
	dal.Reg(b.infoDal, Info{})
	dal.Reg(b.caseDal, Case{})
}

func (b *Bll) RegMsg(msg qf.MessageMap) {

}

func (b *Bll) RefBll() []qf.IBll {
	return nil
}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}

func (b *Bll) SaveInfo(ctx *qf.Context) (interface{}, error) {
	// 获取提交的基本信息
	model := &Info{}
	err := ctx.BindModel(model)
	if err != nil {
		return nil, err
	}
	// 赋值内容
	model.Content = b.BuildContent(model)
	// 保存患者基本信息
	return model, b.infoDal.Save(model)
}

func (b *Bll) SaveCase(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}

func (b *Bll) DeleteCase(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}

func (b *Bll) DeletePatient(ctx *qf.Context) (interface{}, error) {
	return nil, nil
}

//
// GetCase
//  @Description: 通过病历唯一ID获取单条病历信息
//  @param ctx 传入结构体为：
//  @return interface{} 病历Case结构Json
//  @return error
//
func (b *Bll) GetCase(ctx *qf.Context) (interface{}, error) {
	// 外部调用时，请用框架创建上下文
	// qf.BuildContext(map)
	
	m := Case{}
	m.ID = ctx.GetUIntValue("id")
	return b.caseDal.GetModel(m)
}

//
// GetCaseList
//  @Description: 根据病历号或者患者唯一号，获取一组病历列表
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) GetCaseList(ctx *qf.Context) (interface{}, error) {
	// 外部调用时，请用框架创建上下文
	// qf.BuildContext(map)

	caseId := ctx.GetStringValue("caseId")
	infoId := ctx.GetUIntValue("infoId")
	return b.caseDal.Search(infoId, caseId)
}

//
// GetPatientInfo
//  @Description: 根据病历唯一ID获取患者基本信息
//  @param ctx
//  @return interface{} 基本信息Info结构Json
//  @return error
//
func (b *Bll) GetPatientInfo(ctx *qf.Context) (interface{}, error) {
	// 外部调用时，请用框架创建上下文
	// qf.BuildContext(map)

	hisId := ctx.GetStringValue("hisId")
	// 按HIS唯一号查询患者基本信息
	if hisId != "" {
		return b.infoDal.GetInfoByHisId(hisId)
	}
	// 按患者唯一号获取患者基本信息
	m := Info{}
	m.ID = ctx.GetUIntValue("id")
	return b.infoDal.GetModel(m)
}
