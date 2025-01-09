package ttu

import (
	"errors"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"

	"ehang.io/nps/cmd/npc/iping"
)

func PingIpD(ip string) bool {
	log.Println("ping: ip: ", ip)
	cmd := exec.Command("ping", "-c", "1", "-W", "2", ip)
	ret, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("ping: get output error: ", err)
		return false
	}
	// 正则匹配包丢失率
	// X received, X% packet loss
	re := regexp.MustCompile(`(\d+)% packet loss`)
	matched := re.FindStringSubmatch(string(ret))
	if len(matched) < 2 {
		log.Println("ping: can't find packet loss rate")
		return false
	}
	if matched[1] != "100" {
		log.Println("ping: packet loss rate 100%")
		return false
	}
	return true
}

func PingIp(ip, inf string) bool {
	log.Println("ping: ip: ", ip, " net card: ", inf)
	cmd := exec.Command("ping", "-I", inf, "-c", "1", "-W", "2", ip)
	ret, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("ping: get output error: ", err)
		return false
	}
	// 正则匹配包丢失率
	// X received, X% packet loss
	re := regexp.MustCompile(`(\d+)% packet loss`)
	matched := re.FindStringSubmatch(string(ret))
	if len(matched) < 2 {
		log.Println("ping: can't find packet loss rate")
		return false
	}
	if matched[1] == "100" {
		log.Println("ping: packet loss rate 100%", inf)
		return false
	}
	return true
}

func GetAvailableNetCard(ip string) (string, error) {
	card, _ := iping.FindNetInterfaceWhichCanAssessInternet(ip, func(inf net.Interface) bool {
		return strings.Contains(inf.Name, "ppp")
	})
	if card == nil {
		return "", errors.New("no available net card")
	}
	return card.Name, nil
}
