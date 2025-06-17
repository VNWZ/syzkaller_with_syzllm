package syzllm_pkg

import "fmt"

type ServerInfo struct {
	Host     string
	Port     string
	Hostname string
	Env      string
}

var configs = map[string]ServerInfo{
	"test": {
		Host:     "127.0.0.1",
		Port:     "6678",
		Hostname: "parallels",
		Env:      "test",
	},
	"prod": {
		Host:     "127.0.0.1",
		Port:     "6678",
		Hostname: "RIPPLE",
		Env:      "prod",
	},
}

func GetServerInfo(env string) (ServerInfo, error) {
	info, ok := configs[env]
	if !ok {
		return ServerInfo{}, fmt.Errorf("unknown environment: %s", env)
	}
	return info, nil
}
