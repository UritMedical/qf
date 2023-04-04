package main

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/example/bll/patient"
	"github.com/UritMedical/qf/example/bll/sqlserver"
)

func main() {
	qf.Run(regBll, nil)
}

func regBll(s *qf.Service) {
	// 注册自定义业务
	s.RegBll(&patient.Bll{}, "") // 患者信息业务
	s.RegBll(&sqlserver.Bll{}, "")
}
