package main

import (
	"github.com/bitly/go-nsq"
//	"github.com/giskook/smarthome-access/pb"
	"github.com/golang/protobuf/proto"
	"log"
    "github.com/huoyan108/dataRelayGo/olderhc/pbgo"

)

func main() {
	config := nsq.NewConfig()
	var w *nsq.Producer
	w, _ = nsq.NewProducer("192.168.8.178:4150", config)

	replogin := &Report.ControlReport{
		Tid:          189115999977674,
		SerialNumber: 1,
		Command: &Report.Command{
			Type: Report.Command_CMT_REPLOGIN,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: 1,
				},
			},
		},
	}
	reqdata, err := proto.Marshal(replogin)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	log.Println("send topic")

	err = w.Publish("sha2app", reqdata)
	if err != nil {
		log.Panic("Could not connect")
	}

	w.Stop()
}
