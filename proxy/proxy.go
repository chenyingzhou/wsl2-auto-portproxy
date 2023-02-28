package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Proxy struct {
	LocalPort  uint16
	RemotePort uint16
	RemoteIp   string
	Listener   *net.TCPListener
	IsRunning  bool
}

func (p *Proxy) Start() error {
	if p.IsRunning {
		return nil
	}
	localAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", p.LocalPort))
	if err != nil {
		log.Printf("resove local port error, port: %d, err: %v\n", p.LocalPort, err)
		return err
	}
	p.Listener, err = net.ListenTCP("tcp", localAddr)
	if err != nil {
		log.Printf("could not start proxy server on %d: %v\n", p.LocalPort, err)
		return err
	}
	log.Printf("proxy start: %d -> %s:%d", p.LocalPort, p.RemoteIp, p.RemotePort)
	go func() {
		for {
			conn, err := p.Listener.AcceptTCP()
			if err != nil {
				break
			}
			go p.handleTCPConn(conn, 5000)
		}
	}()
	p.IsRunning = true
	return nil
}

func (p *Proxy) Stop() error {
	if !p.IsRunning {
		return nil
	}
	p.IsRunning = false
	log.Printf("proxy stop:  %d -> %s:%d", p.LocalPort, p.RemoteIp, p.RemotePort)
	return p.Listener.Close()
}

func (p *Proxy) handleTCPConn(conn *net.TCPConn, timeout int64) {
	defer func() { _ = conn.Close() }()
	_ = conn.SetKeepAlive(true)
	_ = conn.SetKeepAlivePeriod(time.Second * 15)
	targetAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", p.RemoteIp, p.RemotePort))
	if err != nil {
		log.Printf("resove remote Addr error,%s\n", err)
	}
	c, err := net.DialTimeout("tcp", targetAddr.String(), time.Duration(timeout)*time.Second)
	if err != nil {
		log.Println("Could not connect to remote server:", err)
		return
	}
	client, _ := c.(*net.TCPConn)
	defer func() { _ = client.Close() }()
	log.Printf("client '%v' connected to server '%v'\n", conn.RemoteAddr(), client.RemoteAddr())

	_ = client.SetKeepAlive(true)
	_ = client.SetKeepAlivePeriod(time.Second * 15)

	stop := make(chan bool)

	go func() {
		_, err := io.Copy(client, conn)
		if err != nil {
			log.Printf("client '%v' disconnected from server '%v' by the server\n", conn.RemoteAddr(), client.RemoteAddr())
		}
		stop <- true
	}()

	go func() {
		_, err := io.Copy(conn, client)
		if err != nil {
			log.Printf("client '%v' disconnected from server '%v' by the client\n", conn.RemoteAddr(), client.RemoteAddr())
		}
		stop <- true
	}()

	<-stop
}
