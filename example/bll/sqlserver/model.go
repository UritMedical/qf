package sqlserver

import "time"

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
