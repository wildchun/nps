package main

import (
	"ehang.io/nps/gui/desktop/api"
	"errors"
	"github.com/astaxie/beego/logs"
	"os/exec"
	"strings"
)

const PortStart = 22000
const PortEnd = 23000

func GetESN() (string, error) {
	//devctl -e
	// esn :1205024190010605
	// get esn "1205024190010605" from output
	cmd := exec.Command("devctl", "-e")
	if err := cmd.Run(); err != nil {
		logs.Critical("get esn error: ", err)
		return "", err
	}
	result, err := cmd.Output()
	if err != nil {
		logs.Critical("get esn error: ", err)
		return "", err
	}

	esn := ""
	for _, line := range strings.Split(string(result), "\n") {
		if strings.HasPrefix(line, "esn :") {
			esn = strings.TrimSpace(strings.TrimPrefix(line, "esn :"))
			break
		}
	}
	if esn == "" {
		logs.Critical("get esn error: not found")
		return "", err
	}
	return esn, nil
}

func AutoRegister() error {
	var esn string
	var err error
	if _, err = api.GetKey(); err != nil {
		return err
	}
	if esn, err = GetESN(); err != nil {
		return err
	}
	cltList, err := api.GetList()
	if err != nil {
		return err
	}
	for _, c := range cltList.Rows {
		if c.VerifyKey == esn {
			logs.Info("device already registered")
			return nil
		}
	}
	if err = api.ClientAdd(&api.ClientAddReq{
		Remark:          "TTU:" + esn,
		U:               "",
		P:               "",
		Limit:           0,
		VerifyKey:       esn,
		ConfigConnAllow: false,
		Compress:        false,
		Crypt:           false,
		RateLimit:       "",
		FlowLimit:       "",
		MaxConn:         "",
		MaxTunnel:       "",
	}); err != nil {
		return err
	}

	cltList, err = api.GetList()
	if err != nil {
		return err
	}
	clientId := -1

	for _, c := range cltList.Rows {
		if c.VerifyKey == esn {
			clientId = c.Id
		}
	}

	if clientId == -1 {
		return errors.New("client not found")
	}

	if err = api.ClientAddTcpTunnel(&api.ClientAddTcpTunnelReq{
		Remark:   "LOCAL_SSH",
		ClientId: clientId,
		Target:   "127.0.0.1:8888",
		Port:     PortStart + clientId*2,
	}); err != nil {
		return err
	}
	return nil
}
