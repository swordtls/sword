package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type ProxyServerConfig struct {
	ServerConfig
	RedirectAddress string `json:"redirect_address"`
}


type ServerConfig struct {
	ListenPort string `json:"listen_port"`
	ServerPerm string `json:"server_perm"`
	ServerKey  string `json:"server_key"`
	ClientPerm string `json:"client_perm"`
	CaCert     string `json:"ca_cert"`
}

type ClientConfig struct {
	ListenPort    string `json:"listen_port"`
	ServerAddress string `json:"server_address"`
	ClientPerm    string `json:"client_perm"`
	ClientKey     string `json:"client_key"`
	CaCert        string `json:"ca_cert"`
}

func readConfigFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func ParseProxyServerConfig(filename string) ProxyServerConfig {
	var proxyServerConfig ProxyServerConfig
	configBytes, err := readConfigFile(filename)
	if err != nil {
		panic(fmt.Sprintf("get config file error:%v", err))
	}
	err = json.Unmarshal(configBytes, &proxyServerConfig)
	if err != nil {
		panic(fmt.Sprintf("json decode config file error:%v", err))
	}
	return proxyServerConfig
}

func ParseServerConfig(filename string) ServerConfig {
	var serverConfig ServerConfig
	configBytes, err := readConfigFile(filename)
	if err != nil {
		panic(fmt.Sprintf("get config file error:%v", err))
	}
	err = json.Unmarshal(configBytes, &serverConfig)
	if err != nil {
		panic(fmt.Sprintf("json decode config file error:%v", err))
	}
	return serverConfig
}

func ParseClientConfig(filename string) ClientConfig {
	var clientConfig ClientConfig
	configBytes, err := readConfigFile(filename)
	if err != nil {
		panic(fmt.Sprintf("get config file error:%v", err))
	}
	err = json.Unmarshal(configBytes, &clientConfig)
	if err != nil {
		panic(fmt.Sprintf("json decode config file error:%v", err))
	}
	return clientConfig
}
