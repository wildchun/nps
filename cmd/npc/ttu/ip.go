package ttu

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
)

func parseServerAddrFromExeFileName(exeFileName string) (string, error) {
	// 从文件名解析出服务器地址
	// 例如：npc-124.223.42.242-10011
	// 返回：124.223.42.242:10011
	re := regexp.MustCompile(`npc-(\d+\.\d+\.\d+\.\d+)-(\d+)`)
	matched := re.FindStringSubmatch(exeFileName)
	if len(matched) < 3 {
		return "", errors.New("can't find server address")
	}
	return matched[1] + ":" + matched[2], nil
}

func GetServerAddrFromBinFileName() (string, error) {
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	return parseServerAddrFromExeFileName(exec)
}
