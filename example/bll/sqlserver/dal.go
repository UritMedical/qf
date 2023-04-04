package sqlserver

import (
	"fmt"
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/util/qdate"
	"time"
)

type Dal struct {
	qf.BaseDal
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

//
// getHistory
//  @Description: 获取患者最近6个月的历史结果
//  @param patId
//  @return []Req_Result
//
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
