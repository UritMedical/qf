package main

import (
	"qf"
	"qf/mc/patient"
	"qf/mc/user"
	"qf/util/launcher"
)

func main() {
	launcher.Run(start, stop)
}

var service *qf.Service

func start() {
	service = qf.NewService()
	service.RegBll(&patient.Bll{}, "api")
	service.RegBll(&user.UserBll{}, "api")
	service.Run()
}

func stop() {
	service.Stop()
}
