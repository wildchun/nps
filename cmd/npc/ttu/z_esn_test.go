package ttu

import (
	"regexp"
	"testing"
)

func TestGetEsn(t *testing.T) {
	sec := `
esn :1005024420062941
`
	re := regexp.MustCompile(`esn :(\d+)`)
	matched := re.FindStringSubmatch(sec)
	if len(matched) < 2 {
		t.Fatal("can't find esn")
	}
	if matched[1] != "1005024420062941" {
		t.Fatal("esn not match")
	}
	t.Log("esn: ", matched[1])
}
