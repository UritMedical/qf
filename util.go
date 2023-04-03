package qf

import (
	"github.com/UritMedical/qf/util"
	"github.com/UritMedical/qf/util/launcher"
	"github.com/UritMedical/qf/util/qerror"
	"github.com/UritMedical/qf/util/qreflect"
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

//
// Bind
//  @Description: 将source结构体中的数据绑定到target结构体中
//  @param targetPtr vm结构体, 必须是指针
//  @param source 原始结构体
//  @return error
//
func Bind(targetPtr interface{}, source interface{}) error {
	r := qreflect.New(source)
	return util.SetModel(targetPtr, r.ToMapExpandAll())
}
