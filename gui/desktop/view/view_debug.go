package view

import (
	"ehang.io/nps/gui/desktop/api"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type Debug struct {
	Window fyne.Window
}

func NewDebug() *VClient {
	m := &VClient{}
	m.Window = fyne.CurrentApp().NewWindow("接口测试")
	m.Window.SetContent(m.setupUi())
	return m
}
func (m *Debug) setupUi() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("测试获取密钥", m.onTestGetKeyBtnClicked),
		widget.NewButton("测试获取客户端列表", m.onTestGetClientListBtnClicked),
		widget.NewButton("测试获取客户端", m.onTestGetClientClicked),
	)
}

func (m *Debug) onTestGetKeyBtnClicked() {
	if _, err := api.GetKey(); err != nil {
		dialog.ShowError(err, m.Window)
	} else {
		dialog.ShowInformation("获取成功", "获取成功", m.Window)
	}
}

func (m *Debug) onTestGetClientListBtnClicked() {
	if _, err := api.GetList(); err != nil {
		dialog.ShowError(err, m.Window)
	} else {
		dialog.ShowInformation("获取成功", "获取成功", m.Window)
	}
}

func (m *Debug) onTestGetClientClicked() {
}
