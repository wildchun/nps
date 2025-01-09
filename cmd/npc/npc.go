package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"ehang.io/nps/client"
	"ehang.io/nps/cmd/npc/ttu"
	"ehang.io/nps/lib/common"
	"ehang.io/nps/lib/config"
	"ehang.io/nps/lib/file"
	"ehang.io/nps/lib/install"
	"ehang.io/nps/lib/version"
	"github.com/astaxie/beego/logs"
	"github.com/ccding/go-stun/stun"
	"github.com/kardianos/service"
)

var (
	serverAddr   = flag.String("server", "", "Server addr (ip:port)")
	configPath   = flag.String("config", "", "Configuration file path")
	verifyKey    = flag.String("vkey", "", "Authentication key")
	logType      = flag.String("log", "stdout", "Log output mode（stdout|file）")
	connType     = flag.String("type", "tcp", "Connection type with the server（kcp|tcp）")
	proxyUrl     = flag.String("proxy", "", "proxy socks5 url(eg:socks5://111:222@127.0.0.1:9007)")
	logLevel     = flag.String("log_level", "7", "log level 0~7")
	registerTime = flag.Int("time", 2, "register time long /h")
	localPort    = flag.Int("local_port", 2000, "p2p local port")
	password     = flag.String("password", "", "p2p password flag")
	target       = flag.String("target", "", "p2p target")
	localType    = flag.String("local_type", "p2p", "p2p target")
	logPath      = flag.String("log_path", "", "npc log path")
	debug        = flag.Bool("debug", true, "npc debug")
	pprofAddr    = flag.String("pprof", "", "PProf debug addr (ip:port)")
	stunAddr     = flag.String("stun_addr", "stun.stunprotocol.org:3478",
		"stun server address (eg:stun.stunprotocol.org:3478)")
	ver            = flag.Bool("version", false, "show current version")
	disconnectTime = flag.Int("disconnect_timeout", 60,
		"not receiving check packet times, until timeout will disconnect the client")
)

var NPC_VERSION = "ttu_custom_v1.0.2-20250109"

func main() {
	flag.Parse()
	logs.Reset()
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
	if *ver {
		common.PrintVersion()
		return
	}
	if *logPath == "" {
		*logPath = common.GetNpcLogPath()
	}
	if common.IsWindows() {
		*logPath = strings.Replace(*logPath, "\\", "\\\\", -1)
	}
	if *debug {
		logs.SetLogger(logs.AdapterConsole, `{"level":`+*logLevel+`,"color":true}`)
	} else {
		logs.SetLogger(logs.AdapterFile,
			`{"level":`+*logLevel+`,"filename":"`+*logPath+`","daily":false,"maxlines":100000,"color":true}`)
	}
	logs.Info("npc version:", NPC_VERSION)
	// init service
	options := make(service.KeyValue)
	svcConfig := &service.Config{
		Name:        "Npc",
		DisplayName: "nps内网穿透客户端",
		Description: "custom npc client,server",
		Option:      options,
	}
	if !common.IsWindows() {
		svcConfig.Dependencies = []string{
			"Requires=network.target",
			"After=network-online.target syslog.target",
		}
		svcConfig.Option["SystemdScript"] = install.SystemdScript
		svcConfig.Option["SysvScript"] = install.SysvScript
	}
	for _, v := range os.Args[1:] {
		switch v {
		case "install", "start", "stop", "uninstall", "restart":
			continue
		}
		if !strings.Contains(v, "-service=") && !strings.Contains(v, "-debug=") {
			svcConfig.Arguments = append(svcConfig.Arguments, v)
		}
	}
	svcConfig.Arguments = append(svcConfig.Arguments, "-debug=false")
	prg := &npc{
		exit: make(chan struct{}),
	}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logs.Error(err, "service function disabled")
		run()
		// run without service
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
		return
	}
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "status":
			if len(os.Args) > 2 {
				path := strings.Replace(os.Args[2], "-config=", "", -1)
				client.GetTaskStatus(path)
			}
		case "register":
			flag.CommandLine.Parse(os.Args[2:])
			client.RegisterLocalIp(*serverAddr, *verifyKey, *connType, *proxyUrl, *registerTime)
		case "update":
			install.UpdateNpc()
			return
		case "nat":
			c := stun.NewClient()
			c.SetServerAddr(*stunAddr)
			nat, host, err := c.Discover()
			if err != nil || host == nil {
				logs.Error("get nat type error", err)
				return
			}
			fmt.Printf("nat type: %s \npublic address: %s\n", nat.String(), host.String())
			os.Exit(0)
		case "start", "stop", "restart":
			// support busyBox and sysV, for openWrt
			if service.Platform() == "unix-systemv" {
				logs.Info("unix-systemv service")
				cmd := exec.Command("/etc/init.d/"+svcConfig.Name, os.Args[1])
				err := cmd.Run()
				if err != nil {
					logs.Error(err)
				}
				return
			}
			err := service.Control(s, os.Args[1])
			if err != nil {
				logs.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			return
		case "install":
			service.Control(s, "stop")
			service.Control(s, "uninstall")
			install.InstallNpc()
			err := service.Control(s, os.Args[1])
			if err != nil {
				logs.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			if service.Platform() == "unix-systemv" {
				logs.Info("unix-systemv service")
				confPath := "/etc/init.d/" + svcConfig.Name
				os.Symlink(confPath, "/etc/rc.d/S90"+svcConfig.Name)
				os.Symlink(confPath, "/etc/rc.d/K02"+svcConfig.Name)
			}
			return
		case "uninstall":
			err := service.Control(s, os.Args[1])
			if err != nil {
				logs.Error("Valid actions: %q\n%s", service.ControlAction, err.Error())
			}
			if service.Platform() == "unix-systemv" {
				logs.Info("unix-systemv service")
				os.Remove("/etc/rc.d/S90" + svcConfig.Name)
				os.Remove("/etc/rc.d/K02" + svcConfig.Name)
			}
			return
		}
	}
	s.Run()
}

type npc struct {
	exit chan struct{}
}

func (p *npc) Start(s service.Service) error {
	go p.run()
	return nil
}
func (p *npc) Stop(s service.Service) error {
	close(p.exit)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *npc) run() error {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			logs.Warning("npc: panic serving %v: %v\n%s", err, string(buf))
		}
	}()
	run()
	select {
	case <-p.exit:
		logs.Warning("stop...")
	}
	return nil
}

func run() {
	common.InitPProfFromArg(*pprofAddr)
	// p2p or secret command
	if *password != "" {
		commonConfig := new(config.CommonConfig)
		commonConfig.Server = *serverAddr
		commonConfig.VKey = *verifyKey
		commonConfig.Tp = *connType
		localServer := new(config.LocalServer)
		localServer.Type = *localType
		localServer.Password = *password
		localServer.Target = *target
		localServer.Port = *localPort
		commonConfig.Client = new(file.Client)
		commonConfig.Client.Cnf = new(file.Config)
		go client.StartLocalServer(localServer, commonConfig)
		return
	}
	env := common.GetEnvMap()
	if *serverAddr == "" {
		*serverAddr, _ = env["NPC_SERVER_ADDR"]
	}
	if *verifyKey == "" {
		if err := ttu.ReadyESN(); err != nil {
			logs.Error("read esn error: %s", err)
			return
		}
		*verifyKey = ttu.ESN
		logs.Info("verifyKey use esn: %s", ttu.ESN)
	}
	logs.Info("the version of client is %s, the core version of client is %s", version.VERSION, version.GetVersion())
	linkCardCheck(*serverAddr)
	if *verifyKey != "" && *serverAddr != "" && *configPath == "" {
		go func() {
			failedCnt := 0
			for {
				client.NewRPClient(*serverAddr, *verifyKey, *connType, *proxyUrl, nil, *disconnectTime).Start()
				logs.Info("Client closed! It will be reconnected in five seconds")
				failedCnt++
				if failedCnt > 3 {
					logs.Error("Client closed! It will be reconnected in five seconds")
					os.Exit(0)
				}
				time.Sleep(time.Second * 5)
			}
		}()
	} else {
		logs.Error("serverAddr or verifyKey is empty")
	}
}

func exitLater() {
	os.Exit(0)
}

func linkCardCheck(serverAddr string) {
	ip := strings.Split(serverAddr, ":")[0]
	// 尝试 ping 服务器地址
	if ttu.PingIpD(ip) {
		// 如果 ping 通，直接返回
		return
	}

	pppDev, err := ttu.GetAvailableNetCard(ip)
	if err != nil {
		logs.Error("get ppp net card error: %s", err)
		exitLater()
	}
	logs.Info("get ppp net card success, use %s", pppDev)
	// 获取当前路由
	if ttu.IsHostRouteExist(ip, pppDev) {
		logs.Info("host route exist, just run")
		// 如果路由存在，说明已经是第二次启动了，可以直接运行了
		go func(ip, dev string) {
			quit := make(chan os.Signal, 1)
			// 注册需要关注的信号：SIGINT、SIGTERM、SIGQUIT
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			// 阻塞当前 goroutine 等待信号
			<-quit
			err := ttu.DelHostRoute(ip, dev)
			if err != nil {
				logs.Error("del host route error: %s", err)
			} else {
				logs.Info("del host route success")
			}
		}(ip, pppDev)
		return
	}
	if ttu.AddHostRoute(ip, pppDev) != nil {
		logs.Error("add host route error")
	} else {
		logs.Info("add host route success,wait next start")
	}
	// 如果路由不存在，说明是第一次启动，需要等待一段时间再启动
	exitLater()
}
