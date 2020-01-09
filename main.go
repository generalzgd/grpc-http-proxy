package main

import (
	`os`
	`os/signal`
	`runtime`
	`syscall`

	`github.com/astaxie/beego/logs`

	`grpcproxy/config`
	`grpcproxy/mgr`
)

func init() {
	logger := logs.GetBeeLogger()
	logger.SetLevel(logs.LevelInfo)
	logger.SetLogger(logs.AdapterConsole)
	logger.SetLogger(logs.AdapterFile, `{"filename":"logs/file.log","level":7,"maxlines":1024000000,"maxsize":1024000000,"daily":true,"maxdays":7}`)
	logger.EnableFuncCallDepth(true)
	logger.SetLogFuncCallDepth(3)
	logger.Async(100000)
}

func exit(err error) {
	code := 0
	if err != nil {
		logs.Error("error exit.", err)
		code = 1
	}
	logs.GetBeeLogger().Flush()
	os.Exit(code)
}

func main() {
	var err error
	defer func() {
		exit(err)
	}()

	cfg := config.AppConfig{}
	if err = cfg.Load(); err != nil {
		return
	}

	runtime.GOMAXPROCS(runtime.NumCPU())

	manager := mgr.NewManager()
	manager.Init(cfg)
	if err = manager.ServeHttp(); err != nil {
		return
	}

	// 注册服务
	/*agent := consul.GetConsulRemoteInst(cfg.Consul.Address, cfg.Consul.Token)
	agent.Init(cfg.Basic.Name, "proxy", int(cfg.Proxy.Port), int(cfg.Proxy.Port), cfg.Consul.HealthPort, cfg.Consul.HealthType, cfg.LocalIp)
	if err = agent.Register(nil, nil); err == nil {
		logs.Info("Consul register success.")
	} else {
		logs.Error("Consul register fail.", err)
	}*/

	// catchs system signal
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGTERM)
	sig := <-chSig
	logs.Info("siginal:", sig)

	// 注销服务
	/*if err = agent.Deregister(""); err == nil {
		logs.Info("Consul deregister success.")
	} else {
		logs.Error("Consul deregister fail.", err)
	}
	agent.Destroy()*/
	//manager.Destroy()
}
