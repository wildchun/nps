// Package ttu
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/10 11:29
//	@Update  hd_0411_qxc  2025/1/10 11:29
package ttu

import (
	"github.com/vishvananda/netlink"
)

func NetLinkDaemon() {
	nl, err := netlink.Dial(syscall.NETLINK_ROUTE, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer nl.Close()
}
