package demo

import (
	"qf"
	"qf/helper/content"
)

var (
	RouterReqMain = qf.ApiRouter{Id: "reqMain", Explain: "主信息业务"}
)

type MainBll struct {
	qf.BaseBll
	detailBll *DetailBll
}

func (m *MainBll) RegApis(apis *qf.Apis) {
	apis.Reg(RouterReqMain, "")
}

func (m *MainBll) RegMessages(messages *qf.Messages) {

}

func (m *MainBll) RegReferences(references *qf.References) {
	references.Init(m.detailBll, "引用申请明细业务")
}

func (m *MainBll) Init() (err error) {
	return nil
}

func (m *MainBll) Stop() {

}

func (m *MainBll) BeforeApis(kind qf.EApiKind, content content.Content) (interface{}, error) {
	if kind == qf.EApiKindSubmit {
		// 从内容中获取条码号
		reqInfo := struct{ Barcode string }{}
		err := content.BindJson(&reqInfo)
		if err != nil {
			return nil, err
		}
		// TODO：通过索引表检查条码是否存在
		// ...
	}
	if kind == qf.EApiKindGet {
		// 通过明细业务，获取明细内容
		detail, err := m.detailBll.GetModel(content.ID)
		// 返回主表和明细
		return []interface{}{content, detail}, err
	}
	return nil, nil
}

func (m *MainBll) AfterApis(kind qf.EApiKind, latest []content.Content, old content.Content) (interface{}, error) {
	if kind == qf.EApiKindSubmit {
		// TODO：更新索引
	}
	if kind == qf.EApiKindDelete {
		// TODO：删除索引
	}
	return nil, nil
}
