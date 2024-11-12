package view

import (
	"errors"

	"ehang.io/nps/gui/desktop/api"
	"ehang.io/nps/lib/file"
	"ehang.io/nps/lib/version"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Login struct {
	Window fyne.Window
	ui     struct {
		keyEdit  *widget.Entry
		loginBtn *widget.Button
	}
}

func NewLogin() *Login {
	m := &Login{}
	m.Window = fyne.CurrentApp().NewWindow("登录 " + version.VERSION + " (WildChun)")
	m.Window.SetContent(m.setupUi())
	m.Window.Resize(fyne.NewSize(480, 320))
	return m
}
func (m *Login) setupUi() fyne.CanvasObject {
	m.ui.keyEdit = widget.NewPasswordEntry()
	m.ui.keyEdit.SetText("k4am8wi99s7h2dyl")
	m.ui.loginBtn = widget.NewButton("登录", m.onLoginBtnClicked)
	return container.NewVBox(
		container.New(layout.NewFormLayout(),
			widget.NewLabel("密钥"),
			m.ui.keyEdit,
		),
		m.ui.loginBtn,
	)
}

func (m *Login) Show() {
	m.Window.ShowAndRun()
}

func (m *Login) onLoginBtnClicked() {
	userKey := m.ui.keyEdit.Text
	if userKey == "" {
		dialog.ShowError(errors.New("密钥不能为空"), m.Window)
		return
	}
	m.ui.loginBtn.Text = "登录中..."
	m.ui.loginBtn.Refresh()
	go func(userKey string) {
		defer func() {
			m.ui.loginBtn.Text = "登录"
			m.ui.loginBtn.Refresh()
		}()
		if _, err := api.GetKey(); err != nil {
			dialog.ShowError(errors.New("连不了服务武器:艹"+err.Error()), m.Window)
			return
		}
		cltList, err := api.GetList()
		if err != nil {
			dialog.ShowError(errors.New("鉴权失败:服务器错误"+err.Error()), m.Window)
			return
		}
		for _, c := range cltList.Rows {
			if c.VerifyKey == userKey {
				m.LoginSuccess(c)
				return
			}
		}
		dialog.ShowError(errors.New("鉴权失败:密钥错误"+err.Error()), m.Window)
	}(userKey)
}

func (m *Login) LoginSuccess(c *file.Client) {
	v := NewClient(c)
	v.Window.CenterOnScreen()
	v.Window.Show()
	m.Window.Close()
}
