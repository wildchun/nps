package ttu

import (
	"testing"
)

func TestGetServerAddr(t *testing.T) {
	exe := "./npc-124.223.42.242-10011"
	addr, err := parseServerAddrFromExeFileName(exe)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("server addr: ", addr)
}
