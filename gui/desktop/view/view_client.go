package view

import (
	"strconv"
	"time"

	"ehang.io/nps/client"
	"ehang.io/nps/gui/desktop/api"
	"ehang.io/nps/gui/desktop/inerlog"
	"ehang.io/nps/lib/file"
	"ehang.io/nps/lib/version"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/astaxie/beego/logs"
)

type Client struct {
	Window fyne.Window
	ui     struct {
		tunnelTable *widget.Table
		startBtn    *widget.Button
		logBox      *widget.Entry
	}
	d struct {
		client  *file.Client
		tunnels []*file.Tunnel

		start     bool
		closing   bool
		status    string
		cl        *client.TRPClient
		refreshCh chan struct{}
	}
}

const (
	STATUS_START     = "给我开始穿!"
	STATUS_STOP      = "给我停了!"
	STATUS_RECONNECT = "别急，MD，在重连..."
)

func NewClient(c *file.Client) *Client {
	m := &Client{}
	m.d.client = c
	m.d.status = STATUS_START
	m.d.refreshCh = make(chan struct{})
	m.d.cl = new(client.TRPClient)
	m.Window = fyne.CurrentApp().NewWindow("客户端-内穿" + version.VERSION + " (WildChun)")
	m.Window.SetContent(m.setupUi())
	m.Window.Resize(fyne.NewSize(600, 400))
	return m
}

func (m *Client) setupUi() fyne.CanvasObject {
	info := container.New(
		layout.NewFormLayout(),
		widget.NewLabel("用户"),
		widget.NewLabel(m.d.client.Remark),
	)
	m.ui.tunnelTable = widget.NewTable(
		func() (int, int) {
			return len(m.d.tunnels) + 1, 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			lab := o.(*widget.Label)
			if i.Row == 0 {
				lab.SetText([...]string{"备注", "公网", "本地"}[i.Col])
			} else {
				tunnel := m.d.tunnels[i.Row-1]
				lab.SetText([...]string{tunnel.Remark,
					api.ServerIp + ":" + strconv.Itoa(tunnel.Port),
					tunnel.Target.TargetStr,
				}[i.Col])
			}
		})
	m.ui.tunnelTable.SetColumnWidth(0, 150)
	m.ui.tunnelTable.SetColumnWidth(1, 200)
	m.ui.tunnelTable.SetColumnWidth(2, 200)
	m.ui.tunnelTable.OnSelected = func(id widget.TableCellID) {
		if id.Row == 0 {
			return
		}
		tunnel := m.d.tunnels[id.Row-1]
		switch id.Col {
		case 0:
			m.Window.Clipboard().SetContent(tunnel.Remark)
		case 1:
			m.Window.Clipboard().SetContent(api.ServerIp + ":" + strconv.Itoa(tunnel.Port))
		case 2:
			m.Window.Clipboard().SetContent(tunnel.Target.TargetStr)
		}
	}

	m.onUpdateTunnelBtnClicked()

	m.ui.startBtn = widget.NewButton(m.d.status, m.onStartBtnClicked)
	go func() {
		for {
			<-m.d.refreshCh
			m.ui.startBtn.SetText(m.d.status)
		}
	}()

	updateTunnelBtn := widget.NewButton("更新tunnel", m.onUpdateTunnelBtnClicked)

	m.ui.logBox = widget.NewMultiLineEntry()

	go func() {
		for {
			time.Sleep(time.Second)
			m.ui.logBox.SetText(inerlog.GetLog())
		}
	}()

	return container.NewBorder(
		container.NewHBox(info, updateTunnelBtn, m.ui.startBtn),
		nil,
		nil,
		nil,
		container.NewVSplit(m.ui.tunnelTable, m.ui.logBox),
	)
}

func (m *Client) Show() {
	m.Window.ShowAndRun()
}

func (m *Client) setStatus(status string) {
	m.d.status = status
	m.d.refreshCh <- struct{}{}
}

func (m *Client) onUpdateTunnelBtnClicked() {
	go func() {
		tunnels, err := api.GetTunnel(m.d.client.Id)
		if err != nil {
			return
		}
		m.d.tunnels = tunnels
		m.ui.tunnelTable.Refresh()
	}()
}

func (m *Client) onStartBtnClicked() {
	m.d.start = !m.d.start
	if m.d.start {
		// 启动
		inerlog.Clear()
		m.d.closing = false
		go func() {
			for {
				m.d.cl = client.NewRPClient(api.NpsServer,
					m.d.client.VerifyKey,
					"tcp", "", nil, 60)
				m.setStatus(STATUS_STOP)
				logs.Info("start to connect")
				m.d.cl.Start()
				if m.d.closing {
					logs.Info("client exit")
					return
				}
				logs.Warn("client closed, reconnecting in 5 seconds...")
				m.setStatus(STATUS_RECONNECT)
				time.Sleep(time.Second * 5)
			}
		}()
	} else {
		// 停止
		m.setStatus(STATUS_START)
		m.d.closing = true
		if m.d.cl != nil {
			go m.d.cl.Close()
			m.d.cl = nil
		}
	}
}
