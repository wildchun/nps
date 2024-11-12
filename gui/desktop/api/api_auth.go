package api

import (
	"strconv"
	"time"

	"ehang.io/nps/lib/crypt"
)

const (
	ServerIp  = "124.223.42.242"
	Server    = ServerIp + ":10010"
	NpsServer = ServerIp + ":10011"
	CryptKey  = "wildchunwildchun"
)

var AuthKey string

// map[string][]string

func BuildAuthForm(d map[string][]string) map[string][]string {
	form := make(map[string][]string)
	timestamp := time.Now().Unix()
	form["auth_key"] = []string{crypt.Md5(AuthKey + strconv.FormatInt(timestamp, 10))}
	form["timestamp"] = []string{strconv.FormatInt(timestamp, 10)}
	for k, v := range d {
		form[k] = v
	}
	return form
}
