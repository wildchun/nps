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
	STATUS_START     = "启动"
	STATUS_STOP      = "停止"
	STATUS_RECONNECT = "重连中"
)

func NewClient(c *file.Client) *Client {
	m := &Client{}
	m.d.client = c
	m.d.status = STATUS_START
	m.d.refreshCh = make(chan struct{})
	m.d.cl = new(client.TRPClient)
	m.Window = fyne.CurrentApp().NewWindow("内网穿透客户端 " + version.VERSION + " (WildChun)")
	m.Window.SetContent(m.setupUi())
	m.Window.Resize(fyne.NewSize(480, 320))
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
			return len(m.d.tunnels), 3
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			tunnel := m.d.tunnels[i.Row]
			switch i.Col {
			case 0:
				o.(*widget.Label).SetText(tunnel.Remark)
			case 1:
				o.(*widget.Label).SetText(strconv.Itoa(tunnel.Port))
			case 2:
				o.(*widget.Label).SetText(tunnel.Target.TargetStr)
			}
		})

	go func() {
		tunnels, err := api.GetTunnel(m.d.client.Id)
		if err != nil {
			return
		}
		m.d.tunnels = tunnels
		m.ui.tunnelTable.Refresh()
	}()

	m.ui.startBtn = widget.NewButton(m.d.status, m.onStartBtnClicked)
	go func() {
		for {
			<-m.d.refreshCh
			m.ui.startBtn.SetText(m.d.status)
		}
	}()

	m.ui.logBox = widget.NewMultiLineEntry()

	go func() {
		for {
			time.Sleep(time.Second)
			m.ui.logBox.SetText(inerlog.GetLog())
		}
	}()

	return container.NewBorder(
		container.NewHBox(info, m.ui.startBtn),
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

func (m *Client) onStartBtnClicked() {
	m.d.start = !m.d.start
	if m.d.start {
		// 启动
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
