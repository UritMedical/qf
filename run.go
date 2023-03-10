package qf

import "github.com/UritMedical/qf/util/launcher"

var (
	serv       *Service
	regBllFunc func(s *Service)
	stopFunc   func()
)

//
// Run
//  @Description: 启动
//  @param regBll 注册业务（必须）
//  @param stop 自定义释放
//
func Run(regBll func(s *Service), stop func()) {
	regBllFunc = regBll
	stopFunc = stop
	launcher.Run(doStart, doStop)
}

func doStart() {
	// 创建服务
	serv = newService()
	// 注册业务
	regBllFunc(serv)
	// 启动服务
	serv.run()
}

func doStop() {
	// 执行外部释放
	if stopFunc != nil {
		stopFunc()
	}
	// 停止服务
	serv.stop()
}
