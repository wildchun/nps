package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type ClientAddReq struct {
	Remark          string `json:"remark"`
	U               string `json:"u"`
	P               string `json:"p"`
	Limit           int    `json:"limit"`
	VerifyKey       string `json:"vkey"`
	ConfigConnAllow bool   `json:"config_conn_allow"`
	Compress        bool   `json:"compress"`
	Crypt           bool   `json:"crypt"`
	RateLimit       string `json:"rate_limit"`
	FlowLimit       string `json:"flow_limit"`
	MaxConn         string `json:"max_conn"`
	MaxTunnel       string `json:"max_tunnel"`
}
type StatusResponse struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

func NewClientAddReq(remark, vkey string) *ClientAddReq {
	return &ClientAddReq{
		Remark:          remark,
		U:               "",
		P:               "",
		Limit:           0,
		VerifyKey:       vkey,
		ConfigConnAllow: false,
		Compress:        false,
		Crypt:           false,
		RateLimit:       "",
		FlowLimit:       "",
		MaxConn:         "",
		MaxTunnel:       "",
	}
}

func (r *ClientAddReq) toParamsMap() map[string][]string {
	var params = make(map[string][]string)
	boolStr := func(b bool) string {
		if b {
			return "1"
		} else {
			return "0"
		}
	}
	params["remark"] = []string{r.Remark}
	params["u"] = []string{r.U}
	params["p"] = []string{r.P}
	params["limit"] = []string{strconv.Itoa(r.Limit)}
	params["vkey"] = []string{r.VerifyKey}
	params["config_conn_allow"] = []string{boolStr(r.ConfigConnAllow)}
	params["compress"] = []string{boolStr(r.Compress)}
	params["crypt"] = []string{boolStr(r.Crypt)}
	params["rate_limit"] = []string{r.RateLimit}
	params["flow_limit"] = []string{r.FlowLimit}
	params["max_conn"] = []string{r.MaxConn}
	params["max_tunnel"] = []string{r.MaxTunnel}
	return params
}

func ClientAdd(f *ClientAddReq) error {
	response, err := http.PostForm(GetUrl("/client/add/"),
		BuildAuthForm(f.toParamsMap()))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var resp StatusResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Status != 1 {
		return errors.New(resp.Msg)
	}
	return nil
}

type ClientAddTcpTunnelReq struct {
	Remark   string `json:"remark"`
	ClientId int    `json:"client_id"`
	Target   string `json:"target"`
	Port     int    `json:"port"`
}

func ClientAddTcpTunnel(req *ClientAddTcpTunnelReq) error {
	response, err := http.PostForm(GetUrl("/client/add/"),
		BuildAuthForm(map[string][]string{
			"type":      {"tcp"},
			"remark":    {req.Remark},
			"port":      {strconv.Itoa(req.Port)},
			"target":    {req.Target},
			"client_id": {strconv.Itoa(req.ClientId)},
		}))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var resp StatusResponse
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Status != 1 {
		return errors.New(resp.Msg)
	}
	return nil
}
