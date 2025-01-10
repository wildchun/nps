// Package netlink
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/10 12:16
//	@Update  hd_0411_qxc  2025/1/10 12:16
package main

import (
	"ehang.io/nps/cmd/npc/ttu"
	"github.com/astaxie/beego/logs"
)

type LinkDaemon struct {
	OutDev string
	P2P    *ttu.P2PDaemon
}

func (l *LinkDaemon) Start() {
	l.P2P.OnP2PLinkUpdate = l.onP2PLinkChanged
	if err := l.P2P.Start(); err != nil {
		logs.Error("start p2p daemon error: %v", err)
		return
	}
}

func (l *LinkDaemon) onP2PLinkChanged() {
	p2pInfs := l.P2P.Links()
	if len(p2pInfs) == 0 {
		logs.Info("p2p link has no interface")
	} else if len(p2pInfs) == 1 {
		//  内网卡上线或则外网卡拨号了
		l.onSingleP2PLinkExist()
	} else if len(p2pInfs) == 2 {
		l.onBothP2PLinkExist()
	}
}

func (l *LinkDaemon) onSingleP2PLinkExist() {
	// 内网卡上线或则外网卡拨号了,设置本卡默认路由
	// 删除所有默认路由
	inf := l.P2P.Links()[0]
	err := ttu.DelFullNetRoute(inf.Name)
	logs.Info("del default route: net 0.0.0.0 netmask 0.0.0.0 dev %v, ret : %v", inf.Name, err)
	err = ttu.AddFullNetRoute(inf.Name, "0")
	logs.Info("add default route: net 0.0.0.0 netmask 0.0.0.0 dev %v, ret : %v", inf.Name, err)
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

	_ = ttu.DelFullNetRoute(inner)
	_ = ttu.DelFullNetRoute(outInf)

	_ = ttu.AddFullNetRoute(inner, "0")
	_ = ttu.AddFullNetRoute(outInf, "1")
}

func main() {
}
