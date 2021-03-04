package mqtt

import (
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	base    MQTT.Client
	in, out chan Msg
}

type Msg struct {
	Topic, Payload string
}

type ClientInfo struct {
	ClientId  string
	BrokerUrl string
	Username  string
	Password  string
	KeepAlive int
}

func NewClient(info ClientInfo) (*Client, error) {
	//fmt.Println("Init Mqtt Successfully:", cg)
	client := &Client{}
	client.out = make(chan Msg)
	client.in = make(chan Msg)
	client.base = MQTT.NewClient(MQTT.NewClientOptions().
		AddBroker("tcp://" + info.BrokerUrl).
		SetClientID(info.ClientId).
		SetKeepAlive(time.Duration(info.KeepAlive) * time.Second).
		SetUsername(info.Username).
		SetPassword(info.Password).
		SetDefaultPublishHandler(func(c MQTT.Client, m MQTT.Message) {
			client.out <- Msg{
				Topic:   m.Topic(),
				Payload: string(m.Payload()[:]),
			}
		}).
		SetPingTimeout(1 * time.Second))
	if token := client.base.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	go func() {
		for msg := range client.in {
			token := client.base.Publish(msg.Topic, 0, false, msg.Payload)
			token.Wait()
		}
	}()
	return client, nil
}

func (client *Client) Subscribe(topic string) error {
	token := client.base.Subscribe(topic, 0, nil)
	token.Wait()
	return token.Error()
}

func (client *Client) In() chan Msg {
	return client.in
}

func (client *Client) Out() chan Msg {
	return client.out
}
