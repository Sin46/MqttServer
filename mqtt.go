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

func newClient(info ClientInfo) *Client {
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
		panic(token.Error())
	}

	go func() {
		for msg := range client.in {
			token := client.base.Publish(msg.Topic, 0, false, msg.Payload)
			token.Wait()
		}
	}()
	return client
}

func (client *Client) subscribe(topic string) {
	if token := client.base.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
