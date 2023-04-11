package config

import (
	"encoding/json"
	"fmt"
	"github.com/chenyingzhou/wsl2-tcpproxy/util"
	"io/ioutil"
	"log"
	"os/user"
	"runtime"
	"strconv"
	"strings"
)

var configFileName = ".wsl2-tcpproxy.json"

var configFilePath string

var configFileExample = `{
  "predefined": [
    "6380:6379"
  ],
  "ignore": [
    "443"
  ],
  "custom": [
    "8081:192.168.1.99:8080"
  ]
}`

type Config struct {
	Ignore     []PortIgnore
	Predefined []PortProxy
	Custom     []PortProxy
}

type PortProxy struct {
	LocalPort  uint16
	RemotePort uint16
	RemoteIp   string
}

type PortIgnore uint16

func (pp PortProxy) MarshalJSON() ([]byte, error) {
	if pp.RemoteIp == "" {
		return []byte(fmt.Sprintf("%d:%d", pp.LocalPort, pp.RemotePort)), nil
	} else {
		return []byte(fmt.Sprintf("%d:%s:%d", pp.LocalPort, pp.RemoteIp, pp.RemotePort)), nil
	}
}

func (pp *PortProxy) UnmarshalJSON(data []byte) error {
	var ppStr string
	err := json.Unmarshal(data, &ppStr)
	if err != nil {
		return err
	}
	var localPort64, remotePort64 uint64
	ppParts := strings.Split(ppStr, ":")
	if len(ppParts) == 2 {
		ppParts = []string{ppParts[0], "", ppParts[1]}
	}

	localPort64, err = strconv.ParseUint(ppParts[0], 10, 16)
	if err != nil {
		return err
	}
	pp.LocalPort = uint16(localPort64)
	remotePort64, err = strconv.ParseUint(ppParts[2], 10, 32)
	if err != nil {
		return err
	}
	pp.RemotePort = uint16(remotePort64)
	pp.RemoteIp = ppParts[1]
	return nil
}

func (pi PortIgnore) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%d", pi)), nil
}

func (pi *PortIgnore) UnmarshalJSON(data []byte) error {
	piStr := strings.Trim(string(data), "\" ")
	piNum, err := strconv.ParseUint(piStr, 10, 16)
	if err != nil {
		return err
	}
	*pi = PortIgnore(uint16(piNum))
	return nil
}

func init() {
	userHome, _ := user.Current()
	delimiter := "/"
	if runtime.GOOS == "windows" {
		delimiter = "\\"
	}
	configFilePath = userHome.HomeDir + delimiter + configFileName
	_, err := util.CreateFileIfNotExist(configFilePath, configFileExample)
	if err != nil {
		log.Fatalf("config init error: %s", err)
	}
	log.Println("config file: " + configFilePath)
}

func GetConfig() (Config, error) {
	var out Config
	exists, _ := util.PathExists(configFilePath)
	if !exists {
		return out, nil
	}
	b, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return out, err
	}
	if err = json.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}
