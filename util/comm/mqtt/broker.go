package mqtt

import (
	"github.com/fhmq/hmq/broker"
	"os"
)

// StartBroker 启动Broker服务
func StartBroker(port, WsPort string) error {
	cfg, err := broker.ConfigureConfig(os.Args[1:])
	cfg.Port = port
	cfg.HTTPPort = ""
	cfg.WsPort = WsPort
	cfg.Debug = false
	cfg.WsPath = "/ws"

	b, err := broker.NewBroker(cfg)
	if err != nil {
		return err
	}
	b.Start()
	return nil
}
