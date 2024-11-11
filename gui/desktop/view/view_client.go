package view

import (
	"ehang.io/nps/gui/desktop/api"
	"ehang.io/nps/lib/file"
	"ehang.io/nps/lib/version"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

type VClient struct {
	client  *file.Client
	tunnels []*file.Tunnel
	Window  fyne.Window
	ui      struct {
		tunnelTable *widget.Table
		startBtn    *widget.Button
	}
}

func NewVClient(client *file.Client) *VClient {
	m := &VClient{}
	m.client = client
	m.Window = fyne.CurrentApp().NewWindow("内网穿透客户端 " + version.VERSION + " (WildChun)")
	m.Window.SetContent(m.setupUi())
	return m
}

func (m *VClient) setupUi() fyne.CanvasObject {
	info := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("用户"),
		widget.NewLabel(m.client.Remark),
	)
	m.ui.tunnelTable = widget.NewTable(
		func() (int, int) {
			return len(m.tunnels), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			tunnel := m.tunnels[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(tunnel.Remark)
			case 1:
				o.(*widget.Label).SetText(strconv.Itoa(tunnel.Port))
			case 2:
				o.(*widget.Label).SetText(tunnel.Target.TargetStr)
			}
		})

	m.ui.startBtn = widget.NewButton("启动", m.onLoginBtnClicked)
	go func() {
		tunnels, err := api.GetTunnel(m.client.Id)
		if err != nil {
			return
		}
		m.tunnels = tunnels
		m.ui.tunnelTable.Refresh()
	}()
	return container.NewVBox(info,
		container.NewScroll(m.ui.tunnelTable),
		m.ui.startBtn)
}

func (m *VClient) Show() {
	m.Window.ShowAndRun()
}

func (m *VClient) onLoginBtnClicked() {

}
