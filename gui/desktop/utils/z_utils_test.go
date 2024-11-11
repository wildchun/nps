package utils

import (
	"ehang.io/nps/lib/crypt"
	"encoding/hex"
	"testing"
)

//0607fd7b36c22257dbaee212215fcdd6

func TestAesDecryptByCBC(t *testing.T) {
	//data := []byte("wildchun")
	key := "wildchunwildchun"

	//decrypt := AesDecryptByCBC(encrypted, key)
	//if string(data) != decrypt {
	//	t.Errorf("AesDecryptByCBC() = %v, want %v", decrypt, string(data))
	//}

	encrypted, _ := hex.DecodeString("0607fd7b36c22257dbaee212215fcdd6")
	c, e := crypt.AesDecrypt(encrypted, []byte(key))
	t.Log(string(c), e)
}
