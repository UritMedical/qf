package demo

import (
	"qf"
	"qf/helper/content"
)

var (
	RouterReqDetail = qf.ApiRouter{Id: "reqDetail", Explain: "明细信息业务"}
)

type DetailBll struct {
	qf.BaseBll
}

func (d *DetailBll) RegApis(apis *qf.Apis) {

}

func (d *DetailBll) RegMessages(messages *qf.Messages) {

}

func (d *DetailBll) RegReferences(references *qf.References) {

}

func (d *DetailBll) Init() (err error) {
	return nil
}

func (d *DetailBll) Stop() {

}

func (d *DetailBll) BeforeApis(kind qf.EApiKind, content content.Content) (interface{}, error) {
	return nil, nil
}

func (d *DetailBll) AfterApis(kind qf.EApiKind, latest []content.Content, old content.Content) (interface{}, error) {
	return nil, nil
}
