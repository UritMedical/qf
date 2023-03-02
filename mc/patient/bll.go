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
	// 基本信息
	api.Reg(qf.EKindSave, "", b.SaveInfo)        // 保存患者基本信息
	api.Reg(qf.EKindDelete, "", b.DeletePatient) // 删除患者

	// 病历相关
	api.Reg(qf.EKindSave, "case", b.SaveCase)        // 保存病历
	api.Reg(qf.EKindDelete, "case", b.DeleteCase)    // 删除病历
	api.Reg(qf.EKindGetModel, "case", b.GetCase)     // 通过病历唯一ID获取单条病历信息
	api.Reg(qf.EKindGetList, "cases", b.GetCaseList) // 根据病历号或者患者唯一号，获取一组病历列表

	// 完整信息
	api.Reg(qf.EKindGetModel, "", b.GetPatientFull) // 根据唯一号获取单个患者完整信息（基本+病历列表）
}

func (b *Bll) RegDal(dal qf.DalMap) {
	b.infoDal = &InfoDal{}
	b.caseDal = &CaseBll{}
	dal.Reg(b.infoDal, Patient{})
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
	model := &Patient{}
	err := ctx.Bind(model)
	if err != nil {
		return nil, err
	}
	// 刷新内容
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
	id := ctx.GetUIntValue("id")
	return id, b.infoDal.Delete(id)
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

	//id := ctx.GetUIntValue("id")
	//b.caseDal.GetModel(id, Case{})
	return nil, nil
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
	return nil, nil
	//infoId := ctx.GetUIntValue("infoId")
	//caseId := ctx.GetStringValue("caseId")
	//return b.caseDal.Search(infoId, caseId)
}

//
// GetPatientFull
//  @Description: 根据唯一号获取患者完整信息（基本+病历列表）
//  @param ctx
//  @return interface{}
//  @return error
//
func (b *Bll) GetPatientFull(ctx *qf.Context) (interface{}, error) {
	pkg := struct {
		Patient         // 患者基本信息
		CaseList []Case // 包含的病历列表
	}{}
	// 按患者唯一号，获取患者基本信息
	id := ctx.GetUIntValue("id")
	err := b.infoDal.GetModel(id, &pkg)
	if pkg.Id == 0 {
		return nil, nil
	}
	// 根据患者唯一号，获取所有病历列表
	return pkg, err
}
