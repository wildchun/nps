package main

import (
	"ehang.io/nps/gui/desktop/api"
	"ehang.io/nps/lib/file"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

const ServerIp = "124.223.42.242"

func main() {
	if _, err := api.GetKey(); err != nil {
		return
	}
	cltList, err := api.GetList()
	if err != nil {
		return
	}

	ttuFiles := make([]*file.Client, 0)
	for _, c := range cltList.Rows {
		if !strings.HasPrefix(c.Remark, "新疆TTU:") {
			continue
		}
		ttuFiles = append(ttuFiles, c)
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetActiveSheet(index)
	f.SetCellValue("Sheet1", "A1", "备注")
	f.SetCellValue("Sheet1", "B1", "V-KEY")
	f.SetCellValue("Sheet1", "C1", "通道")
	f.SetCellValue("Sheet1", "D1", "公网映射")
	f.SetCellValue("Sheet1", "E1", "本地映射")
	rowIndex := 1

	for _, ttuFile := range ttuFiles {
		tunnels, err := api.GetTunnel(ttuFile.Id)
		if err != nil {
			continue
		}
		for _, tunnel := range tunnels {
			f.SetCellValue("Sheet1", "A"+fmt.Sprint(rowIndex+1), ttuFile.Remark)
			f.SetCellValue("Sheet1", "B"+fmt.Sprint(rowIndex+1), ttuFile.VerifyKey)
			f.SetCellValue("Sheet1", "C"+fmt.Sprint(rowIndex+1), tunnel.Remark)
			f.SetCellValue("Sheet1", "D"+fmt.Sprint(rowIndex+1), fmt.Sprintf("%s:%d", ServerIp, tunnel.Port))
			f.SetCellValue("Sheet1", "E"+fmt.Sprint(rowIndex+1), tunnel.Target.TargetStr)
		}
		rowIndex++
	}
	if err := f.SaveAs("TTU.xlsx"); err != nil {
		fmt.Println(err)
	}
}
