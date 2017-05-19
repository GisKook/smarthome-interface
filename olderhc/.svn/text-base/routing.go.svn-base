package olderhc

import (
	"fmt"
	"os"
    "log"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

var appUpChan = make(chan string)
var c *MQTT.Client

func MsgRecvfun(client *MQTT.Client, msg MQTT.Message) {
	appUpChan <- msg.Topic()
	fmt.Printf("appUp")
	fmt.Printf("topic:[%s]  ", msg.Topic())
	fmt.Printf("Mesg:%s\n", msg.Payload())

}

func SendMsg(pubTipc string,msg string){

    log.Println("mqtt Send:"+msg)
	if token := c.Publish(pubTipc, 0, false,msg); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
func AppCommUnitInit() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883").SetClientID("router-sample")
	opts.SetCleanSession(true)

	c = MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
    log.Println("Connected to MqttServer")
    

	if token := c.Subscribe("appUp", 0,MsgRecvfun ); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

}
func Disconnect(time uint){
	c.Disconnect(time)
}
