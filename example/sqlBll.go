package main

import (
	"fmt"
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/util/qdate"
	"strings"
	"time"
)

type Bll struct {
	qf.BaseBll

	dal *Dal
}

type Dal struct {
	qf.BaseDal
}

type Req_Sample struct {
	Req_Id        string
	Test_Date     time.Time
	Test_Group    string
	Test_Sample   string
	Pat_Id        string
	Pat_DH        string
	Pat_Type      string
	Req_Time      time.Time
	Req_TestTime  time.Time
	Req_CheckTime time.Time
	Req_Info      string
	Sample_Type   string
}

type Req_Detail struct {
	Req_Sample
	Group_Id string
	Gp_Name  string
	Gp_Type  string
}

type Req_Result struct {
	Req_Detail
	Item_InId   string
	Item_Id     string
	Item_Name   string
	Item_Limit  string
	Item_Unit   string
	Test_Result string
}

func (b *Bll) RegApi(a qf.ApiMap) {
	a.Reg(qf.EApiKindGetList, "labs", b.getList)
}

func (b *Bll) RegDal(d qf.DalMap) {
	b.dal = &Dal{}
	d.Reg(b.dal, nil)
}

//
// getResults
//  @Description: 获取条件内的所有样本结果列表
//  @return []Req_Result
//
func (dal *Dal) getResults(ds time.Time, de time.Time, dh string, st string, dp string) []Req_Result {
	tx := dal.DB().
		Table("Req_Sample1 a").
		Select("a.*, b.Group_Id, b.Gp_Name, c.Item_InId, c.Item_Id, c.Item_Name, c.Item_Limit, c.Item_Unit, e.DictValue as Gp_Type").
		Joins("left join Req_Detail b on a.Req_Id = b.Req_Id").
		Joins("left join Req_Result c on a.Test_Date = c.Test_Date and a.Test_Group = c.Test_Group and a.Test_SampleNo = c.Test_SampleNo").
		Joins("left join Lis_Group d on d.Group_Id = b.Group_Id").
		Joins("left join Com_Dict e on d.Group_Belong = e.DictID and e.DictType = 'TP'").
		Where("Req_Time >= ? and Req_Time <= ?", ds, de)
	if dh != "" {
		tx = tx.Where("Pat_DH = ?", dh)
	}
	if st != "" {
		tx = tx.Where("Pat_Type = ?", st)
	}
	if dp != "" {
		tx = tx.Where("Req_DPId = ?", dp)
	}
	// 查询
	list := make([]Req_Result, 0)
	err := tx.Find(&list).Error
	if err != nil {
		fmt.Println(err)
	}
	return list
}

func (dal *Dal) getHistory(patId string) []Req_Result {
	ds, _ := qdate.Parse(qdate.ToString(time.Now(), "yyyy-MM-dd"), "yyyy-MM-dd")
	tx := dal.DB().
		Table("Req_Sample a").
		Select("a.*, b.Item_InId, b.Item_Id, b.Item_Name, b.Item_Limit, b.Item_Unit").
		Joins("left join Req_Result b on a.Test_Date = b.Test_Date and a.Test_Group = b.Test_Group and a.Test_SampleNo = b.Test_SampleNo").
		Where("a.Pat_Id = ? and a.Test_Date >= ?", patId, ds.AddDate(0, -6, 0))
	// 查询
	list := make([]Req_Result, 0)
	err := tx.Find(&list).Error
	if err != nil {
		fmt.Println(err)
	}
	return list
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
