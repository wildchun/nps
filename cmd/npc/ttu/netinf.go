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
	"github.com/siddontang/go/log"
)

type NetInf struct {
}

func CreateNetInf() *NetInf {
	o := &NetInf{}
	return o
}

func (n *NetInf) FindAvailableNetCard(host string) (net.Interface, error) {
	p := iping.PingOption{
		Count:      2,
		Timeout:    2000,
		Size:       32,
		Nerverstop: false,
	}
	infs, _ := net.Interfaces()

	for _, iface := range infs {
		if !strings.Contains(iface.Name, "ppp") {
			continue
		}
		err := AddHostRoute(host, iface.Name)
		if err != nil {
			return net.Interface{}, err
		}
		log.Info("try interface: ", iface.Name)
		if ip, err := iping.GetInfAddress(&iface); err == nil {
			if p.Ping3(host, ip) {
				return iface, nil
			}
		}
		DelHostRoute(host, iface.Name)
	}
	return net.Interface{}, errors.New("no interface can access " + host)
}

func (n *NetInf) HasP2PNetCard() bool {

	return false
}
