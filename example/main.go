package main

import (
	"github.com/UritMedical/qf"
	"github.com/UritMedical/qf/mc/patient"
	"github.com/UritMedical/qf/user"
)

func main() {
	qf.Run(regBll, nil)
}

func regBll(s *qf.Service) {
	// 注册框架提供的通用业务
	// 通用业务位于mc文件夹内
	s.RegBll(&user.Bll{}, "")    // 用户业务
	s.RegBll(&patient.Bll{}, "") // 患者信息业务
	// 注册自定义业务
	// ...
}
