package mgr

import (
	`context`
	`errors`
	`fmt`
	`net/http`
	`strings`
	`time`

	`github.com/astaxie/beego/logs`
	libs `github.com/generalzgd/comm-libs`
	grpclb_consul `github.com/generalzgd/grpc-svr-frame/grpc-consul`
	ctrl `github.com/generalzgd/grpc-svr-frame/grpc-ctrl`
	`github.com/golang/protobuf/proto`
	`github.com/grpc-ecosystem/grpc-gateway/runtime`
	grpcpool `github.com/processout/grpc-go-pool`
	`google.golang.org/grpc`

	`grpcproxy/config`
	"grpcproxy/define"
	`grpcproxy/iproto`
)

type Manager struct {
	ctrl.GrpcController
	cfg            config.AppConfig
	httpInfoCenter *define.ClientInfoCenter
}

func NewManager() *Manager {
	return &Manager{
		GrpcController: ctrl.MakeGrpcController(),
		httpInfoCenter: define.NewClientInfoCenter(),
	}
}

func (p *Manager) Init(cfg config.AppConfig) {
	p.cfg = cfg

	if p.cfg.LocalIp == "" {
		p.cfg.LocalIp = libs.GetInnerIp()
	}
	//prewarn.SetSendMailCallback(func(body string) {
	// 发送邮件
	//title := fmt.Sprintf("%s - %s Warning(%s)", p.cfg.Basic.Name, libs.GetEnvName(), p.cfg.LocalIp)
	//mail.SendMailByUrl(title, body, cfg.MailAddr...)
	//})
	//statistic.SetWarnHandler(prewarn.NewWarn)
	grpclb_consul.InitRegister(p.cfg.Consul.Address)
}

//
func (p *Manager) ServeHttp() error {
	addr := fmt.Sprintf(":%d", p.cfg.Proxy.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// todo 目标服务是否启用tls

	// 测试代码
	// p.httpInfoCenter.Put("12F44847-85EB-E6D8-6D38-09AA118BB528", &common.ClientConnInfo{
	// 	Uid:      163,
	// 	Platform: 1,
	// 	State:    true,
	// })

	err := iproto.RegisterRestfulGatewayHandlerClient(ctx, mux, opts, p.getEndpointAddrByMeth, p.getGrpcClientConnByMeth,
		p.httpCallBeginHandler, p.httpCallDoneHandler, p.httpCallQps)
	if err != nil {
		logs.Error("serve http gate fail.", err)
		return err
	}

	go func() {
		http.ListenAndServe(addr, mux)
	}()
	// p.svrMap[cfg.Name] = server
	logs.Info("start serve %s(%s) with secure(%v)", p.cfg.Name, addr, p.cfg.Proxy.Secure)
	return nil
}

// @param meth string: package.TargetService/Method
func (p *Manager) getEndpointAddrByMeth(meth string) string {
	_, tarSvr, _, _ := iproto.ParseMethod(meth)
	tarSvr = strings.ToLower(tarSvr + "svr")
	cfg, ok := p.cfg.EndpointSvr[tarSvr]
	if ok {
		return cfg.Address
	}
	return ""
}

func (p *Manager) getGrpcClientConnByMeth(meth string) (*grpcpool.ClientConn, error) {
	_, tarSvr, _, _ := iproto.ParseMethod(meth)
	tarSvr = strings.ToLower(tarSvr + "svr")
	if cfg, ok := p.cfg.EndpointSvr[tarSvr]; ok {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		return p.GetGrpcConnWithLB(cfg, ctx)
	}
	return nil, errors.New("endpoint config empty:" + tarSvr)
}

// grpc转发前的回调处理
func (p *Manager) httpCallBeginHandler(meth string, req *http.Request) (int, bool) {
	//statistic.NewRecord(statistic.Stat_Tps)

	// todo 如果有cookie,则生成对应的clientinfo

	/*tmp, _ := dcopy.InstanceToMap(info)
	for k, v := range tmp {
		// header key需要有MetadataHeaderPrefix开头，才会识别出来并转发给endpoint; 终端取到的key不含MetadataHeaderPrefix
		req.Header.Add(runtime.MetadataHeaderPrefix+k, libs.Interface2String(v)) // MetadataHeaderPrefix
	}*/

	return http.StatusAccepted, true
}

// grpc结束回调处理
func (p *Manager) httpCallDoneHandler(meth string, reply proto.Message, w http.ResponseWriter, req *http.Request) {

}

func (p *Manager) httpCallQps(d time.Duration) {
	//statistic.NewRecord(statistic.Stat_Qps, d)
}
