package olderhc

import (
	"encoding/json"
	nsq "github.com/bitly/go-nsq"
	"github.com/giskook/smarthome-interface/base"
	"github.com/giskook/smarthome-interface/olderhc/pbgo"
	"log"
	"strconv"
)

const ADD_DEVICE string = "1"
const DEL_DEVICE string = "2"
const NOTIFICATION string = "3"
const NOTIFY_ONOFF string = "4"
const NOTIFY_LEVEL string = "5"

type AddDevice struct {
	DeviceID   string
	DeviceName string
	Endpoints  []EndpointTag
}

type ProtocolAddDevice struct {
	Protocol string
	Content  *AddDevice
}

func _parse_add_device(paras []*Report.Command_Param) *ProtocolAddDevice {
	count := int(paras[3].Npara)
	endpoints := make([]EndpointTag, count)
	pos := 4
	for i := 0; i < count; i++ {
		endpoints[i].Endpoint = uint8(paras[pos].Npara)
		pos++
		endpoints[i].Devicetype = uint16(paras[pos].Npara)
		pos++
		if endpoints[i].Devicetype == 0x0402 {
			endpoints[i].Zonetype = uint16(paras[pos].Npara)
			pos++
		} else {
			endpoints[i].Zonetype = 0
		}
	}

	return &ProtocolAddDevice{
		Protocol: ADD_DEVICE,
		Content: &AddDevice{
			DeviceID:   base.Uint2HexString(uint64(paras[0].Npara)),
			DeviceName: string(paras[1].Strpara),
			Endpoints:  endpoints,
		},
	}
}

type DelDevice struct {
	DeviceID string
}

type ProtocolDelDevice struct {
	Protocol string
	Content  *DelDevice
}

func _parse_del_device(paras []*Report.Command_Param) *ProtocolDelDevice {
	return &ProtocolDelDevice{
		Protocol: DEL_DEVICE,
		Content: &DelDevice{
			DeviceID: base.Uint2HexString(paras[0].Npara),
		},
	}
}

type notification struct {
	Deviceid   string
	Endpoint   uint8
	Warntime   uint64
	Zonetype   uint16
	Zonestatus uint16
}

type protocol_notification struct {
	Protocol string
	Content  *notification
}

func _parse_notification(paras []*Report.Command_Param) *protocol_notification {
	return &protocol_notification{
		Protocol: NOTIFICATION,
		Content: &notification{
			Deviceid:   base.Uint2HexString(paras[0].Npara),
			Endpoint:   uint8(paras[1].Npara),
			Warntime:   uint64(paras[2].Npara),
			Zonetype:   uint16(paras[3].Npara),
			Zonestatus: uint16(paras[4].Npara),
		},
	}
}

type notify_onoff struct {
	Deviceid string
	Endpoint uint8
	Status   uint8
}

type protocol_notify_onoff struct {
	Protocol string
	Content  *notify_onoff
}

func _parse_notify_onoff(paras []*Report.Command_Param) *protocol_notify_onoff {
	return &protocol_notify_onoff{
		Protocol: NOTIFY_ONOFF,
		Content: &notify_onoff{
			Deviceid: base.Uint2HexString(paras[1].Npara),
			Endpoint: uint8(paras[0].Npara),
			Status:   uint8(paras[2].Npara),
		},
	}
}

type notify_level struct {
	Deviceid string
	Endpoint uint8
	Status   uint8
}

type protocol_notify_level struct {
	Protocol string
	Content  *notify_level
}

func _parse_notify_level(paras []*Report.Command_Param) *protocol_notify_level {
	return &protocol_notify_level{
		Protocol: NOTIFY_LEVEL,
		Content: &notify_level{
			Deviceid: base.Uint2HexString(paras[1].Npara),
			Endpoint: uint8(paras[2].Npara),
			Status:   uint8(paras[3].Npara),
		},
	}
}

func Receive(cs map[string]chan []*Report.Command_Param) {
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		defer func() {
			if x := recover(); x != nil {
				log.Println("sorry, chan has closed")
			}
		}()

		//log.Printf("Got one message: %v", message.Body)

		res := BuildResult(message.Body)
		if res.commandType == int32(Report.Command_CMT_REP_ADD_DEL_DEVICE) {

			var paras []*Report.Command_Param
			paras = []*Report.Command_Param(res.paras)

			var jsonRes []byte //返回内容

			action := int(paras[2].Npara)
			if action == 1 {
				protocol_add_dev := _parse_add_device(paras)
				log.Println(paras)
				jsonRes, _ = json.Marshal(protocol_add_dev)
			} else {
				protocol_del_dev := _parse_del_device(paras)
				jsonRes, _ = json.Marshal(protocol_del_dev)
			}

			tmpmac, _ := strconv.Atoi(res.mac)
			topic := strconv.FormatInt(int64(tmpmac), 16)
			SendMsg(topic, string(jsonRes))
		} else if res.commandType == int32(Report.Command_CMT_REP_NOTIFICATION) {

			var jsonRes []byte //返回内容
			var paras []*Report.Command_Param
			paras = []*Report.Command_Param(res.paras)

			protocol_notify := _parse_notification(paras)

			//mapValue := make(map[string]string)
			//mapValue["deviceid"] = strconv.Itoa(int(paras[0].Npara))
			//mapValue["ednpoint"] = strconv.Itoa(int(paras[1].Npara))
			//mapValue["warntime"] = strconv.Itoa(int(paras[2].Npara))
			//mapValue["zonetype"] = strconv.Itoa(int(paras[3].Npara))
			//mapValue["zonestatus"] = strconv.Itoa(int(paras[4].Npara))
			//jsonRes, _ = json.Marshal(mapValue)
			jsonRes, _ = json.Marshal(protocol_notify)
			log.Println(jsonRes)
			tmpmac, _ := strconv.Atoi(res.mac)
			topic := strconv.FormatInt(int64(tmpmac), 16)
			SendMsg(topic, string(jsonRes))
		} else if res.commandType == int32(Report.Command_CMT_REP_NOTIFY_ONOFF) {
			var jsonRes []byte //返回内容
			var paras []*Report.Command_Param
			paras = []*Report.Command_Param(res.paras)

			protocol_notify_onoff := _parse_notify_onoff(paras)
			jsonRes, _ = json.Marshal(protocol_notify_onoff)
			log.Println(jsonRes)
			tmpmac, _ := strconv.Atoi(res.mac)
			topic := strconv.FormatInt(int64(tmpmac), 16)
			log.Println("send notify onoff")
			SendMsg(topic, string(jsonRes))
		} else if res.commandType == int32(Report.Command_CMT_REP_NOTIFY_LEVEL) {
			var jsonRes []byte //返回内容
			var paras []*Report.Command_Param
			paras = []*Report.Command_Param(res.paras)

			protocol_notify_level := _parse_notify_level(paras)
			jsonRes, _ = json.Marshal(protocol_notify_level)
			log.Println(jsonRes)
			tmpmac, _ := strconv.Atoi(res.mac)
			topic := strconv.FormatInt(int64(tmpmac), 16)
			SendMsg(topic, string(jsonRes))
		} else {
			log.Println("recv nsq mac " + res.mac)
			log.Println("recv nsq ser " + strconv.Itoa(int(res.serialNum)))
			chanKey := res.mac + strconv.Itoa(int(res.serialNum))
			log.Println("recv nsq " + chanKey)
			ci, ok := cs[chanKey]
			log.Println(cs)
			if !ok {
				log.Println("no chan!")
			} else {
				ci <- res.paras
				close(ci)
				delete(cs, chanKey)
				//log.Println("delete existed chan!")
			}
		}
		return nil
	}))

	nsqlookup_url, _ := config.GetString("message", "nsqlookup_url")
	errConsumer := q.ConnectToNSQD(nsqlookup_url)
	if errConsumer != nil {
		log.Panic("Consumer could not connect nsq")
	}
}

func Send(command []byte) {

	command_topic, _ := config.GetString("message", "command_topic")

	err := w.Publish(command_topic, command)
	if err != nil {
		log.Panic("communicator Send could not connect nsq")
	}

	//log.Println("send command ok!", command)
}
