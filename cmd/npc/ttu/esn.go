package ttu

import (
	"github.com/pkg/errors"
	"os/exec"
	"regexp"
)

/**
sysadm@SCT230A:~$ devctl -e
esn :1005024420062941
*/

var ESN = ""

func ReadyESN() error {
	cmd := exec.Command("devctl", "-e")
	ret, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	re := regexp.MustCompile(`esn :(\d+)`)
	matched := re.FindStringSubmatch(string(ret))
	if len(matched) < 2 {
		return errors.New("can't find esn")
	}
	ESN = matched[1]
	return nil
}
