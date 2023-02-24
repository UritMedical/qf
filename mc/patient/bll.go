package patient

import "qf"

type Bll struct {
	qf.BaseBll
	baseDal *BaseDal
}

func (b *Bll) RegApi(api qf.ApiMap) {
	api.Reg(qf.EKindSave, "base", b.saveBase)
}

func (b *Bll) RegDal(dal qf.DalMap) {
	b.baseDal = &BaseDal{}
	dal.Reg(b.baseDal, Base{})
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

func (b *Bll) saveBase(ctx *qf.Context) (interface{}, error) {
	// 获取提交的基本信息
	model := &Base{}
	err := ctx.BindModel(model)
	if err != nil {
		return nil, err
	}
	// 赋值内容
	model.Content = b.BuildContent(model)
	// 保存患者基本信息
	return model, b.baseDal.Save(model)
}
