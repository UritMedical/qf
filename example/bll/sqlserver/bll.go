package sqlserver

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/util/qdate"
	"strings"
)

type Bll struct {
	qf.BaseBll

	dal *Dal
}

func (b *Bll) RegApi(a qf.ApiMap) {
	a.Reg(qf.EApiKindGetList, "labs", b.getList)
}

func (b *Bll) RegDal(d qf.DalMap) {
	b.dal = &Dal{}
	d.Reg(b.dal, nil)
}

func (b *Bll) RegFault(f qf.FaultMap) {

}

func (b *Bll) RegMsg(m qf.MessageMap) {

}

func (b *Bll) RegRef(r qf.RefMap) {

}

func (b *Bll) Init() error {
	return nil
}

func (b *Bll) Stop() {

}

func (b *Bll) getList(ctx *qf.Context) (interface{}, qf.IError) {
	// 查询条件
	ds, e1 := qdate.Parse(ctx.GetStringValue("ds"), "yyyy-MM-dd")
	de, e2 := qdate.Parse(ctx.GetStringValue("de"), "yyyy-MM-dd")
	if e1 != nil || e2 != nil {
		return nil, qf.Error(9999, "自己的故障码")
	}
	dh := ctx.GetStringValue("dh")
	st := ctx.GetStringValue("st")
	dp := ctx.GetStringValue("dp")

	// 查询
	list := b.dal.getResults(ds, de, dh, st, dp)

	// 获取列表
	pats := map[string][]Req_Result{}
	for _, r := range list {
		pats[r.Pat_Id] = nil
	}
	// 获取所有患者6个月前的样本
	for k, _ := range pats {
		pats[k] = b.dal.getHistory(k)
	}

	// 转换为待输出的结构
	for _, r := range list {
		// 解析INFO
		infos := transInfo(r.Req_Info)
		// 提取科室名称和医生名字
		_ = infos["Pat_DocName"]
		_ = infos["Pat_DPName"]
		// 然后比较获取历史结构
	}

	return list, nil
}

func transInfo(info string) map[string]string {
	final := map[string]string{}

	for _, s := range strings.Split(info, ",") {
		sp := strings.Split(s, "=")
		final[sp[0]] = sp[1]
	}

	return final
}
