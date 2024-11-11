package api

import (
	"ehang.io/nps/lib/crypt"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
)

type GetKeyResponse struct {
	CryptAuthKey string `json:"crypt_auth_key"`
	CryptType    string `json:"crypt_type"`
	Status       int    `json:"status"`
}

func GetKey() (string, error) {
	if AuthKey != "" {
		return AuthKey, nil
	}
	response, err := http.Post(GetUrl("/auth/getauthkey"), "application/json", nil)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)
	var getKeyResponse GetKeyResponse
	err = json.NewDecoder(response.Body).Decode(&getKeyResponse)
	if err != nil {
		return "", err
	}

	encrypted, err := hex.DecodeString(getKeyResponse.CryptAuthKey)
	if err != nil {
		return "", err
	}
	if auth, err := crypt.AesDecrypt(encrypted, []byte(CryptKey)); err != nil {
		return "", err
	} else {
		AuthKey = string(auth)
		return AuthKey, nil
	}
}
