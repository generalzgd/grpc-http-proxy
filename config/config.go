package config

import (
	`os`
	`path/filepath`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/grpc-svr-frame/config/ymlcfg`
)

// 网关连接配置
type ProxyConfig struct {
	Secure      bool              `yaml:"secure"`    // false: http/tcp/ws  true: https/tls/wss, 版本号默认1.1
	CertFiles   []ymlcfg.CertFile `yaml:"certfiles"` // 证书文件，pem格式
	BufferSize  int               `yaml:"buffersize"`
	MaxConn     int               `yaml:"maxconn"`
	IdleTimeout time.Duration     `yaml:"idletimeout"`
	Port        uint32            `yaml:"port"` // 侦听端口
}

type AppConfig struct {
	Name        string                           `yaml:"name"`
	Ver         string                           `yaml:"ver"`
	Memo        string                           `yaml:"memo"`
	Proxy       ProxyConfig                      `yaml:"proxy"` // http代理端口
	Consul      ymlcfg.ConsulConfig              `yaml:"consul"`
	EndpointSvr map[string]ymlcfg.EndpointConfig `yaml:"endpoint"` // 要转发的终端
	MailAddr    []string                         `yaml:"mailaddr"`
	LocalIp     string                           `yaml:"localip"`
}

func (p *AppConfig) Load() error {
	path := filepath.Join(filepath.Dir(os.Args[0]), "config", "config_dev.yml")
	logs.Info("load config: %s", path)

	return ymlcfg.LoadYaml(path, p)
}
