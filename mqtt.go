package mqtt

import (
	"fmt"
	"os"
	"time"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type mqttInfo struct {
	CLIENT_ID string
    BROKER_URL string
	USERNAME string
	PASSWORD string
	KEEPALIVE int
}

type config struct {
    Mqtt mqttInfo
}

type MqttSub struct{
    Enable bool
    Topic string
}
type MqttMsg struct{
    Topic,  Payload string
}

//整成一个协程,通过 chan 来
func Sub(c MQTT.Client, sub_ch chan MqttSub) {
	for {
		x := <-sub_ch
		fmt.Println("Sub")
		if x.Enable {
			if token := c.Subscribe(x.Topic, 0, nil); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}			
		}else {
			if token := c.Unsubscribe(x.Topic); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}	
		}

	}
}

func Server(f MQTT.MessageHandler, pub_ch chan MqttMsg, sub_ch chan MqttSub){
	var cg config
	cg.Mqtt = mqttInfo{
		CLIENT_ID : "MAC0111",
		BROKER_URL : "61.131.1.193:1883",
		USERNAME : "FjdzMacUser",
		PASSWORD :"geomacuser",
		KEEPALIVE : 2,
	}

	fmt.Println("Init Mqtt Successfully:",cg.Mqtt)
	opts := MQTT.NewClientOptions().AddBroker("tcp://"+cg.Mqtt.BROKER_URL).SetClientID(cg.Mqtt.CLIENT_ID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetUsername(cg.Mqtt.USERNAME)
	opts.SetPassword(cg.Mqtt.PASSWORD)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go Sub(c, sub_ch)
	for {
		x := <-pub_ch
		//有多个生产者，但只有一个消费者
		token := c.Publish(x.Topic, 0, false, x.Payload)
		token.Wait()
	}
}

//声明MQTT的接口
var mqtt_pub_ch = make(chan MqttMsg)
var mqtt_sub_ch = make(chan MqttSub)
var mqtt_sub_func MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
    // CONTROLLER.SubHandler(mqtt_pub_ch, msg.Topic(), msg.Payload())
	fmt.Println("get message")
}

// func main(){
// 	var temp MqttSub

// 	go Server(mqtt_sub_func, mqtt_pub_ch, mqtt_sub_ch)
// 	temp.Enable=true
// 	temp.Topic="Publish_MAC009"
// 	mqtt_sub_ch <- temp
// 	temp.Topic="Subscription_MAC009"
// 	mqtt_sub_ch <- temp
// 	for {
//         //redis某个键值更新 : pub_ch <- sendMsg
//     }
// }