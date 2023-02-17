package mqtt

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type msg struct {
	Topic   string
	Qos     byte
	Retain  bool
	payload interface{}
}

type Client struct {
	client    mqtt.Client
	sendChan  chan msg
	stopChan  chan struct{}
	topicDict map[string]mqtt.MessageHandler
}

//NewClient host例如："tcp://iot.eclipse.org:1883"
func NewClient(host, clientId string) *Client {

	pc := &Client{

		sendChan:  make(chan msg),
		stopChan:  make(chan struct{}),
		topicDict: make(map[string]mqtt.MessageHandler),
	}
	opts := mqtt.NewClientOptions().AddBroker(host).SetClientID(clientId).SetConnectRetry(true).SetConnectTimeout(time.Second * 5)
	opts.SetOnConnectHandler(pc.doSubscribe)
	c := mqtt.NewClient(opts)
	pc.client = c

	pc.client.Connect()
	go pc.run()

	return pc
}

func (p *Client) doSubscribe(client mqtt.Client) {
	for k, v := range p.topicDict {
		client.Unsubscribe(k)
		client.Subscribe(k, 1, v)
	}
}

func (p *Client) run() {
	for true {
		select {
		case <-p.stopChan:
			return
		case m := <-p.sendChan:
			t := p.client.Publish(m.Topic, m.Qos, m.Retain, m.payload)
			_ = t.WaitTimeout(time.Second * 10)
		}
	}
}

//Publish 发布消息
func (p *Client) Publish(topic string, qos byte, payload interface{}, retain bool) {
	payloadJson, _ := json.Marshal(payload)
	m := msg{
		Topic:   topic,
		Qos:     qos,
		payload: payloadJson,
		Retain:  retain,
	}
	if p.client.IsConnectionOpen() {
		p.sendChan <- m
	}

}
func (p *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) {
	p.topicDict[topic] = callback
	if p.client.IsConnectionOpen() {
		p.client.Subscribe(topic, qos, callback)
	}

}
func (p *Client) Unsubscribe(topic string) {
	p.client.Unsubscribe(topic)
}

//Close 关闭客户端
func (p *Client) Close() {
	close(p.stopChan)
	if p.client != nil {
		p.client.Disconnect(250)
	}
}
