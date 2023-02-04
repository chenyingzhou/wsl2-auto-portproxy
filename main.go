package main

import (
	"github.com/chenyingzhou/wsl2-tcpproxy/config"
	"github.com/chenyingzhou/wsl2-tcpproxy/proxy"
	"github.com/chenyingzhou/wsl2-tcpproxy/service"
	"github.com/chenyingzhou/wsl2-tcpproxy/storage"
	"log"
	"time"
)

func main() {
	// get config interval
	go func() {
		for {
			c, err := config.GetConfig()
			if err != nil {
				log.Printf("error getting config file: %s", err)
			} else {
				storage.Conf = c
			}
			time.Sleep(time.Second * 5)
		}
	}()
	for {
		// get linux's ip
		storage.WslIp, _ = service.GetWslIP()
		// get all tcp ports in linux
		storage.WslPorts = service.GetWslPorts()
		// get all tcp ports in windows
		storage.WinPorts = service.GetWinPorts()
		// ignore ports
		for _, ignore := range storage.Conf.Ignore {
			for i, port := range storage.WslPorts {
				if port == uint16(ignore) {
					storage.WslPorts = append(storage.WslPorts[:i], storage.WslPorts[i+1:]...)
				}
			}
		}
		// merge wsl oldProxy and custom oldProxy
		newProxyPool := make(map[uint16]*proxy.Proxy)
		for _, remotePort := range storage.WslPorts {
			localPort := remotePort
			for _, item := range storage.Conf.Predefined {
				if remotePort == item.RemotePort {
					localPort = item.LocalPort
					break
				}
			}
			newProxyPool[localPort] = &proxy.Proxy{
				LocalPort:  localPort,
				RemotePort: remotePort,
				RemoteIp:   storage.WslIp,
			}
		}
		for _, item := range storage.Conf.Custom {
			newProxyPool[item.LocalPort] = &proxy.Proxy{
				LocalPort:  item.LocalPort,
				RemotePort: item.RemotePort,
				RemoteIp:   item.RemoteIp,
			}
		}
		// migrate oldProxy and stop outdated oldProxy
		for localPort, oldProxy := range storage.ProxyPool {
			newProxy, ok := newProxyPool[localPort]
			if ok && newProxy.RemotePort == oldProxy.RemotePort && newProxy.RemoteIp == oldProxy.RemoteIp {
				newProxyPool[localPort] = oldProxy
			} else {
				_ = oldProxy.Stop()
			}
		}
		// start new proxy
		for _, newProxy := range newProxyPool {
			omitted := false
			for _, winPort := range storage.WinPorts {
				if winPort == newProxy.LocalPort {
					omitted = true
					break
				}
			}
			if !omitted {
				_ = newProxy.Start()
			}
		}
		storage.ProxyPool = newProxyPool
		time.Sleep(time.Second * 5)
	}
}
