package olderhc

import (
	nsq "github.com/bitly/go-nsq"
	"github.com/msbranco/goconfig"
	"log"
    "github.com/giskook/smarthome-interface/olderhc/pbgo"
	"os"
)

//指令接收者
var q *nsq.Consumer

//指令发送者
var w *nsq.Producer

//配置文件
var config *goconfig.ConfigFile

//主机名
var hostname string

// 监听端口

var Port string

//channel 存储器，key为command serialNum
//var cs map[string]chan byte = make(map[string]chan byte)
  var cs map[string] chan []*Report.Command_Param = make(map[string]chan []*Report.Command_Param)
var delay int64

func init() {

	//初始化配置
	var configError error
	config, configError = goconfig.ReadConfigFile("conf.ini")
	if configError != nil {
		log.Panic("can not read the config file")
	}

	//初始化主机名
	var errHost error
	hostname, errHost = os.Hostname()
	if errHost != nil {
		log.Fatal("Hostname can not obtain!")
	}
    AppCommUnitInit()

	//端口
	Port, _ = config.GetString("server", "port")

	//初始化指令接收者
	configConsumer := nsq.NewConfig()
    consumerTopic,_ := config.GetString("message","consumerTopic")
	q, _ = nsq.NewConsumer(consumerTopic, "ch", configConsumer)

	//初始化指令发送者
	configProducer := nsq.NewConfig()
	nsq_url, _ := config.GetString("message", "nsq_url")
	w, _ = nsq.NewProducer(nsq_url, configProducer)

	delay, _ = config.GetInt64("message", "delay_time")

	//指令回复处理handler
	Receive(cs)
}
