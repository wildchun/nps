package api

import (
	"ehang.io/nps/lib/file"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type GetListResponse struct {
	BridgePort int            `json:"bridgePort"`
	BridgeType string         `json:"bridgeType"`
	Ip         string         `json:"ip"`
	Rows       []*file.Client `json:"rows"`
	Total      int            `json:"total"`
}

func GetList() (*GetListResponse, error) {
	response, err := http.PostForm(GetUrl("/client/list/"),
		BuildAuthForm(url.Values{
			"start": []string{"0"},
			"limit": []string{"100"},
		}))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var getListResponse GetListResponse
	err = json.NewDecoder(response.Body).Decode(&getListResponse)
	if err != nil {
		return nil, err
	}
	return &getListResponse, nil
}
