package olderhc

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/smarthome-interface/base"
	"github.com/giskook/smarthome-interface/olderhc/pbgo"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"time"
)

type EndpointTag struct {
	Endpoint   uint8
	Devicetype uint16
	Zonetype   uint16
	Status     uint8
}
type Device struct {
	Deviceid   string
	Devicename string
	Online     uint8
	Endpoints  []EndpointTag
}
type GatewayInfo struct {
	Gatewayname string
	TotalCount  int
	Devices     []Device
}

func GetZBNodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GetZBNodeHandler")

	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix())
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			log.Println(string(jsonRes))
			fmt.Fprint(w, string(jsonRes))
			log.Println(x)
		}
	}()

	sGatewayid := r.Form["gatewayid"][0]
	tid := base.Macaddr2uint64(sGatewayid)
	//devType := r.Form["devtype"][0]
	//nDevType, _ := strconv.Atoi(devType)
	// 请求参数不正确
	if sGatewayid == "" || len(sGatewayid) != 12 {
		//if sGatewayid == "" {
		jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
		log.Println(string(jsonRes))
		fmt.Fprint(w, string(jsonRes))
		return
	}
	//构造指令内容
	req := &Report.ControlReport{
		Tid:          tid,
		SerialNumber: uint32(s),
		Command: &Report.Command{
			Type:  Report.Command_CMT_REQ_DEVICELIST,
			Paras: []*Report.Command_Param{
			//	&Report.Command_Param{
			//		Type:  Report.Command_Param_UINT64,
			//		Npara: uint64(nDevType),
			//	},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)

	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	ci := GetHttpRouter().SendRequest(chanKey)
	log.Println("send msg")
	Send(reqdata)

	select {
	case res := <-ci:
		log.Println(res)
		value := []*Report.Command_Param(res)[0].Npara
		//log.Println(value)
		if value == NOT_ONLINE {
			jsonRes, _ = json.Marshal(map[string]uint64{"result": value})
		} else {

			var paras []*Report.Command_Param
			paras = res

			var gatewayinfo GatewayInfo
			gatewayinfo.Gatewayname = paras[1].Strpara
			var totalCount int
			totalCount = int(paras[2].Npara)
			log.Println("totalCount:", totalCount)
			gatewayinfo.TotalCount = totalCount
			devicepos := 3
			for i := 0; i < totalCount; i++ {

				var device Device
				//	device.Deviceid = uint64(paras[devicepos+0].Npara)
				//var x int64
				x := int64(paras[devicepos+0].Npara)
				b_buf := bytes.NewBuffer([]byte{})
				binary.Write(b_buf, binary.BigEndian, x)
				fmt.Println("deviceid byte:", b_buf.Bytes(), "int:", x)
				device.Deviceid = base.Uint2Deviceid(b_buf.Bytes())
				//sDeviceid := strconv.Itoa(int(paras[devicepos+0].Npara))
				//cDeviceid := []byte(sDeviceid)
				//device.Deviceid = base.Uint2Deviceid(cDeviceid)
				device.Devicename = paras[devicepos+1].Strpara
				log.Println("device:", device)
				device.Online = uint8(paras[devicepos+2].Npara)
				var count int
				count = int(paras[devicepos+3].Npara)
				log.Println("count:", count)
				devicepos = devicepos + 2 + 2
				//endpointpos := devicepos
				var pos = 2
				for j := 0; j < count; j++ {
					var endpoint EndpointTag
					//endpoint.Endpoint = uint8(paras[endpointpos+j*pos+0].Npara)
					endpoint.Endpoint = uint8(paras[devicepos+0].Npara)
					endpoint.Devicetype = uint16(paras[devicepos+1].Npara)
					if endpoint.Devicetype == 0x0402 {
						endpoint.Zonetype = uint16(paras[devicepos+2].Npara)
						pos = 3
					} else if endpoint.Devicetype == 0x0009 || endpoint.Devicetype == 0x0200 || endpoint.Devicetype == 0x0002 {
						endpoint.Status = uint8(paras[devicepos+2].Npara)
						pos = 3
					} else {
						pos = 2
					}
					devicepos = devicepos + pos
					log.Println("endpoint:", endpoint)
					device.Endpoints = append(device.Endpoints, endpoint)
				}

				gatewayinfo.Devices = append(gatewayinfo.Devices, device)
			}
			log.Println(gatewayinfo)
			jsonRes, _ = json.Marshal(gatewayinfo)
		}
	case <-time.After(time.Duration(delay) * time.Second):
		log.Println("res : 超时")
		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))

	GetHttpRouter().DelRequest(chanKey)
}
