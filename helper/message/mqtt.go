package message

import (
	"encoding/json"
	"fmt"
	"qf/util/comm/mqtt"
)

func NewMessageByMQTT(ip, port string, wsPort string) *ByMQTT {
	_ = mqtt.StartBroker(port, wsPort)
	m := &ByMQTT{
		client: mqtt.NewClient(fmt.Sprintf("tcp://%s:%s", ip, port), "ByMQTT"),
	}
	return m
}

type ByMQTT struct {
	client *mqtt.Client
}

func (m ByMQTT) Publish(msgId string, payload interface{}) error {
	list := [2]string{}
	list[0] = msgId
	p, _ := json.Marshal(payload)
	list[1] = string(p)
	m.client.Publish("DiscoverTrigger", 0, list, true)
	return nil
}
