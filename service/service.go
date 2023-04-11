package service

import (
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Port struct {
	Type      string
	Port      int64
	ProxyPort int64
}

var wslValid bool

func init() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("wsl", "--exec", "pwd")
		_, err := cmd.Output()
		if err == nil {
			wslValid = true
		}
	}
	if !wslValid {
		log.Printf("wsl is not available\n")
	}
}

func GetWslIP() string {
	if !wslValid {
		return ""
	}
	cmd := exec.Command("wsl", "--", "bash", "-c", "ip -4 a show eth0 | grep -oP '(?<=inet\\s)\\d+(\\.\\d+){3}'")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	ip := strings.Replace(string(output), "\n", "", -1)
	reg := regexp.MustCompile("^\\d{1,3}.\\d{1,3}.\\d{1,3}.\\d{1,3}$")
	if !reg.MatchString(ip) {
		return ""
	}
	return ip
}
func GetWslPorts() []uint16 {
	if !wslValid {
		return nil
	}
	cmd := exec.Command("wsl", "--exec", "netstat", "-tunlp")
	reg := regexp.MustCompile("(tcp)\\s+\\d+\\s+\\d+\\s+(:::|0.0.0.0:)(\\d{2,5})")
	return getPorts(cmd, reg, 3)
}

func GetWinPorts() []uint16 {
	if !wslValid {
		return nil
	}
	cmd := exec.Command("cmd", "/c", "Netstat", "-ano", "|", "findstr", "LISTENING")
	reg := regexp.MustCompile("(TCP)\\s+(\\[::\\]:|0.0.0.0:)(\\d{2,5})")
	return getPorts(cmd, reg, 3)
}

func getPorts(cmd *exec.Cmd, reg *regexp.Regexp, position int) []uint16 {
	var ports []uint16
	output, err := cmd.Output()
	if err != nil {
		log.Printf("failed to exec cmd: %s\n", strings.Join(cmd.Args, " "))
		return ports
	}
	rets := reg.FindAllStringSubmatch(string(output), -1)
	for _, ret := range rets {
		duplicated := false
		port, _ := strconv.ParseInt(ret[position], 10, 16)
		for _, find := range ports {
			if find == uint16(port) {
				duplicated = true
				break
			}
		}
		if !duplicated {
			ports = append(ports, uint16(port))
		}
	}
	return ports
}
