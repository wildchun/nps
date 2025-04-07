package main

import (
	_ "ehang.io/nps/gui/desktop/inerlog"
	"ehang.io/nps/gui/desktop/view"
	"ehang.io/nps/lib/common"
	"ehang.io/nps/lib/daemon"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/astaxie/beego/logs"
	"github.com/flopp/go-findfont"
	"os"
	"strings"
)

func init() {
	// 设置中文字体:解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttf") ||
			strings.Contains(path, "simhei.ttf") ||
			strings.Contains(path, "simsun.ttc") ||
			strings.Contains(path, "simkai.ttf") {
			_ = os.Setenv("FYNE_FONT", path)
			break
		}
	}
}

func main() {
	daemon.InitDaemon("npc", common.GetRunPath(), common.GetTmpPath())
	_ = logs.SetLogger("file", `{"filename":"npc_desktop/npc.log"}`)
	fyne.SetCurrentApp(app.New())
	w := view.NewLogin().Window
	w.CenterOnScreen()
	w.ShowAndRun()
}
