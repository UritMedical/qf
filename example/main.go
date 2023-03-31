package main

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/patient"
)

func main() {
	qf.Run(regBll, nil)
}

func regBll(s *qf.Service) {
	// 注册自定义业务
	s.RegBll(&patient.Bll{}, "") // 患者信息业务
	s.RegBll(&Bll{}, "")
}
