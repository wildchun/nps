package api

import (
	"ehang.io/nps/lib/file"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type GetTunnelResponse struct {
	Rows  []*file.Tunnel `json:"rows"`
	Total int            `json:"total"`
}

func GetTunnel(clientId int) ([]*file.Tunnel, error) {
	response, err := http.PostForm(GetUrl("/index/gettunnel/"),
		BuildAuthForm(url.Values{
			"type":      []string{"tcp"},
			"offset":    []string{"0"},
			"limit":     []string{"100"},
			"client_id": []string{strconv.Itoa(clientId)},
		}))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var getTunnelResponse GetTunnelResponse
	err = json.NewDecoder(response.Body).Decode(&getTunnelResponse)
	if err != nil {
		return nil, err
	}
	return getTunnelResponse.Rows, nil
}
