// Package ttu
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/10 15:50
//	@Update  hd_0411_qxc  2025/1/10 15:50
package ttu

import (
	"os"
	"time"

	"github.com/astaxie/beego/logs"
)

type LinkDaemon struct {
	OutDev string
	P2P    *P2PDaemon
	d      func(f func())
}

func (l *LinkDaemon) Start() {
	l.d = NewDebounce(time.Millisecond * 500)
	l.P2P.OnP2PLinkUpdate = l.onP2PLinkChanged
	if err := l.P2P.Start(); err != nil {
		logs.Error("start p2p daemon error: %v", err)
		return
	}
	logs.Info("start p2p daemon ,out dev: %v", l.OutDev)
}

func (l *LinkDaemon) onP2PLinkChanged() {
	l.d(func() {
		p2pInfs := l.P2P.Links()
		if len(p2pInfs) == 0 {
			logs.Info("p2p link has no interface")
			os.Exit(0)
		} else if len(p2pInfs) == 1 {
			//  内网卡上线或则外网卡拨号了
			logs.Info("p2p link has one interface")
			l.onSingleP2PLinkExist()
		} else if len(p2pInfs) == 2 {
			logs.Info("p2p link has two interface")
			l.onBothP2PLinkExist()
		}
	})
}

func (l *LinkDaemon) onSingleP2PLinkExist() {
	// 内网卡上线或则外网卡拨号了,设置本卡默认路由
	// 删除所有默认路由
	inf := l.P2P.Links()[0]
	err := DelFullNetRoute(inf.Name)
	logs.Info("del default route: net 0.0.0.0 netmask 0.0.0.0 dev %v, ret : %v", inf.Name, err)
	err = AddFullNetRoute(inf.Name, "0")
	logs.Info("add default route: net 0.0.0.0 netmask 0.0.0.0 dev %v, ret : %v", inf.Name, err)

	// 如果外网卡掉线了  退出程序
	if inf.Name == l.OutDev {
		logs.Error("out dev %v offline, exit", l.OutDev)
		os.Exit(0)
	}
}

func (l *LinkDaemon) onBothP2PLinkExist() {
	// 内网卡上线和外网卡拨号了
	// 找出那个是内网卡,那个是外网卡
	infs := l.P2P.Links()
	inner := infs[0].Name
	if inner == l.OutDev {
		inner = infs[1].Name
	}
	outInf := l.OutDev

	_ = DelFullNetRoute(inner)
	_ = DelFullNetRoute(outInf)

	_ = AddFullNetRoute(inner, "0")
	_ = AddFullNetRoute(outInf, "1")
}
