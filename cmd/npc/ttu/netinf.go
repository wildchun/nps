// Package ttu
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/9 16:08
//	@Update  hd_0411_qxc  2025/1/9 16:08
package ttu

import (
	"errors"
	"net"
	"strings"

	"ehang.io/nps/cmd/npc/iping"
	"github.com/astaxie/beego/logs"
	errors2 "github.com/pkg/errors"
)

type NetInf struct {
}

func CreateNetInf() *NetInf {
	o := &NetInf{}
	return o
}

func (n *NetInf) GetP2PNetCard() []net.Interface {
	infs, _ := net.Interfaces()
	var p2pInfs []net.Interface
	for _, iface := range infs {
		if strings.Contains(iface.Name, "ppp") {
			p2pInfs = append(p2pInfs, iface)
		}
	}
	return p2pInfs
}

func (n *NetInf) FindAvailableNetCard(host string, infs []net.Interface) (net.Interface, error) {
	p := iping.PingOption{
		Count:      2,
		Timeout:    2000,
		Size:       32,
		Nerverstop: false,
	}
	for _, iface := range infs {
		if !strings.Contains(iface.Name, "ppp") {
			continue
		}
		err := AddHostRoute(host, iface.Name)
		if err != nil {
			return net.Interface{}, errors2.Wrap(err, "add host route error")
		}
		logs.Info("try interface: %v", iface.Name)
		if ip, err := iping.GetInfAddress(&iface); err == nil {
			if p.Ping3(host, ip) {
				return iface, nil
			}
		}
		_ = DelHostRoute(host, iface.Name)
	}
	return net.Interface{}, errors.New("no interface can access " + host)
}

func (n *NetInf) HasP2PNetCard() bool {

	return false
}
