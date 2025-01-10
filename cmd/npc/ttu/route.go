package ttu

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Item struct {
	DstIp   string
	DstMask string
	ViaIp   string
	DevNic  string
}

/*
route record file : /proc/net/route
Iface       Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
ens38   0000FEA9        00000000        0001    0       0       1000    0000FFFF        0       0       0
ens38   0000A8C0        00000000        0001    0       0       0       00FFFFFF        0       0       0
ens38   0000A8C0        00000000        0001    0       0       101     00FFFFFF        0       0       0
ens33   0002A8C0        00000000        0001    0       0       0       00FFFFFF        0       0       0
ens33   0002A8C0        00000000        0001    0       0       100     00FFFFFF        0       0       0
*/

func TransIpStr(ipStr string) string {
	ip := make([]byte, 4)
	for i := 6; i >= 0; i -= 2 {
		num := byte(0)
		if ipStr[i] >= '0' && ipStr[i] <= '9' {
			num += (ipStr[i] - '0') * 16
		} else if ipStr[i] >= 'A' && ipStr[i] <= 'Z' {
			num += (ipStr[i] - 'A' + 10) * 16
		} else {
			num += (ipStr[i] - 'a' + 10) * 16
		}
		if ipStr[i+1] >= '0' && ipStr[i+1] <= '9' {
			num += ipStr[i+1] - '0'
		} else if ipStr[i+1] >= 'A' && ipStr[i+1] <= 'Z' {
			num += ipStr[i+1] - 'A' + 10
		} else {
			num += ipStr[i+1] - 'a' + 10
		}
		ip[i/2] = num
	}
	return strconv.Itoa(int(ip[3])) + "." + strconv.Itoa(int(ip[2])) + "." + strconv.Itoa(int(ip[1])) + "." + strconv.Itoa(int(ip[0]))
}

func ReadSysRoute() ([]Item, error) {
	file, err := os.Open("/proc/net/route")
	row := 11
	if err != nil {
		return nil, err
	}
	defer file.Close()
	tmp, err := ioutil.ReadAll(file)
	routeStrs := strings.Fields(string(tmp))
	length := len(routeStrs) / row
	var routes []Item
	for i := 1; i < length; i++ {
		index := i * row
		route := Item{
			DstIp:   TransIpStr(routeStrs[index+1]),
			DstMask: TransIpStr(routeStrs[index+7]),
			ViaIp:   TransIpStr(routeStrs[index+2]),
			DevNic:  routeStrs[index],
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func AddHostRoute(dstIp string, dev string) error {
	items, err := ReadSysRoute()
	if err != nil {
		return err
	}
	for _, item := range items {
		if item.DstIp == dstIp && item.DevNic == dev {
			return nil
		}
		if item.DstIp == dstIp {
			// 删除原来的路由
			_ = DelHostRoute(dstIp, item.DevNic)
		}
	}
	cmd := exec.Command("route", "add", "-host", dstIp, "dev", dev)
	return cmd.Run()
}
func DelHostRoute(dstIp string, dev string) error {
	items, err := ReadSysRoute()
	if err != nil {
		return errors.Wrap(err, "read sys route error")
	}
	exist := false
	for _, item := range items {
		if item.DstIp == dstIp && item.DevNic == dev {
			exist = true
			break
		}
	}
	if !exist {
		return nil
	}
	cmd := exec.Command("route", "del", "-host", dstIp, "dev", dev)
	return cmd.Run()
}

func IsHostRouteExist(dstIp string, dev string) bool {
	items, err := ReadSysRoute()
	if err != nil {
		return false
	}
	for _, item := range items {
		if item.DstIp == dstIp && item.DevNic == dev {
			return true
		}
	}
	return false
}

func AddNetRoute(net, netmask, metric, dev string) error {
	cmd := exec.Command("route", "add", "-net", net, "netmask", netmask, "metric", metric, "dev", dev)
	return cmd.Run()
}

func DelNetRoute(net, netmask, dev string) error {
	cmd := exec.Command("route", "del", "-net", net, "netmask", netmask, "dev", dev)
	return cmd.Run()
}
