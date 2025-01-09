// Package iping
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/9 13:22
//	@Update  hd_0411_qxc  2025/1/9 13:22
package iping

import (
	"testing"
)

func TestPing(t *testing.T) {
	// p := NewPingOption()
	// p.ping3("www.baidu.com", nil, net.Interface{})
}

func TestPing2(t *testing.T) {
	inf, err := FindNetInterfaceWhichCanAssessInternet("www.baidu.com")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(inf)
}
