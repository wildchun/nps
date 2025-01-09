// Package iping
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/9 13:44
//	@Update  hd_0411_qxc  2025/1/9 13:44
package iping

import (
	"fmt"
	"net"
)

func FindNetInterfaceWhichCanAssessInternet(host string, filter ...func(inf net.Interface) bool) (*net.Interface,
	error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	filterFunc := func(inf net.Interface) bool {
		return true
	}

	if len(filter) > 0 {
		filterFunc = filter[0]
	}

	p := PingOption{
		Count:      1,
		Timeout:    2000,
		Size:       32,
		Nerverstop: false,
	}
	for _, iface := range interfaces {
		if !filterFunc(iface) {
			continue
		}
		fmt.Println("try interface: ", iface.Name)
		if ip, err := GetInfAddress(&iface); err != nil {
			continue
		} else {
			if p.Ping3(host, ip) {

				return &iface, nil
			}
		}
	}
	return nil, fmt.Errorf("no interface can access %s", host)
}
