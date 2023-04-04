package qf

import (
	"github.com/UritMedical/qf/util/launcher"
	"github.com/UritMedical/qf/util/qerror"
)

//
// Run
//  @Description: 启动
//  @param regBll 注册业务（必须）
//  @param stop 自定义释放
//
func Run(regBll func(s *Service), stop func()) {
	// 收集异常
	defer qerror.Recover(nil)

	regBllFunc = regBll
	stopFunc = stop
	launcher.Run(doStart, doStop)
}
