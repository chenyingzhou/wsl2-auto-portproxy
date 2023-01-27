package storage

import (
	"github.com/chenyingzhou/wsl2-tcpproxy/config"
	"github.com/chenyingzhou/wsl2-tcpproxy/proxy"
)

var ProxyPool map[uint16]*proxy.Proxy

var WslIp string

var WslPorts []uint16

var WinPorts []uint16

var Conf config.Config
