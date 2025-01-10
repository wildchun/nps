// Package ttu
//
//	@Description
//	@Author  hd_0411_qxc  2025/1/10 11:29
//	@Update  hd_0411_qxc  2025/1/10 11:29
package ttu

import (
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
)

type OnP2PLinkUpdate func()

type P2PLink struct {
	Name      string
	OperState netlink.LinkOperState
}
type P2PDaemon struct {
	linkMu          sync.Mutex
	links           map[string]P2PLink
	OnP2PLinkUpdate OnP2PLinkUpdate
	subDone         chan struct{}
}

func NewP2PDaemon() *P2PDaemon {
	return &P2PDaemon{
		links: make(map[string]P2PLink),
	}
}

func (p *P2PDaemon) Links() []P2PLink {
	p.linkMu.Lock()
	defer p.linkMu.Unlock()
	links := make([]P2PLink, 0, len(p.links))
	for _, link := range p.links {
		links = append(links, link)
	}
	return links
}

func (p *P2PDaemon) linkUpdate() {
	infs, err := netlink.LinkList()
	if err != nil {
		logs.Error("netlink.LinkList error: %v", err)
		return
	}

	p.linkMu.Lock()
	defer p.linkMu.Unlock()
	links := make(map[string]P2PLink)
	for _, inf := range infs {
		if strings.Contains(inf.Attrs().Name, "ppp-") {
			links[inf.Attrs().Name] = P2PLink{
				Name:      inf.Attrs().Name,
				OperState: inf.Attrs().OperState,
			}
		}
	}
	hasChanged := len(p.links) == len(links)
	p.links = links
	if !hasChanged && p.OnP2PLinkUpdate != nil {
		p.OnP2PLinkUpdate()
	}
}

func (p *P2PDaemon) Start() error {
	p.linkUpdate()
	ch := make(chan netlink.LinkUpdate)
	p.subDone = make(chan struct{})
	if err := netlink.LinkSubscribe(ch, p.subDone); err != nil {
		return errors.Wrap(err, "netlink.LinkSubscribe")
	}
	go func() {
		for _ = range ch {
			p.linkUpdate()
		}
	}()
	return nil
}

func (p *P2PDaemon) Stop() {
	close(p.subDone)
}
