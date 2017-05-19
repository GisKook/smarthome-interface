package olderhc

import (
	"encoding/json"
	"fmt"
	"github.com/giskook/smarthome-interface/base"
	"github.com/giskook/smarthome-interface/olderhc/pbgo"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"strconv"
	"time"
)

//开关设备指令处理
/*
CMT_REQ_ONOFF
  参数 1.deviceid(uint64) 2.endpoint(uint8) 3.action(uint8 0 for on 1 for off)
     CMT_REP_ONOFF
	   参数 1.result(uint8 0.success 1.fail)
*/
func OnOffDevHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("OnOffDevHandler++++")
	r.ParseForm()
	//now := time.Now()
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			log.Println(string(jsonRes))
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!")
		}
	}()

	sGatewayid := r.Form["gatewayid"][0]
	tid := base.Macaddr2uint64(sGatewayid)
	sDeviceid := r.Form["deviceid"][0]
	deviceid := base.Deviceid2uint64(sDeviceid)
	sEndpoint := r.Form["endpoint"][0]
	nEndpoint, _ := strconv.Atoi(sEndpoint)
	nAction, _ := strconv.Atoi(r.Form["action"][0])
	s, _ := strconv.ParseUint(r.Form["serialid"][0], 16, 32)
	//s = GenerateKey(tid, s)

	// 请求参数不正确
	if sGatewayid == "" || len(sGatewayid) != 12 || sDeviceid == "" || len(sDeviceid) != 16 || sEndpoint == "" {
		jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
		log.Println(string(jsonRes))
		fmt.Fprint(w, string(jsonRes))
		return
	}
	//构造指令内容
	req := &Report.ControlReport{
		Tid:          uint64(tid),
		SerialNumber: uint32(s),
		Command: &Report.Command{
			Type: Report.Command_CMT_REQ_ONOFF,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT64,
					Npara: deviceid,
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nEndpoint),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nAction),
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)

	//log.Printf("control command request: ", reqdata)

	//chanKey := loginid + strconv.Itoa(int(s))
	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	//log.Println("http"+chanKey)
	ci := GetHttpRouter().SendRequest(chanKey)
	Send(reqdata)

	select {
	case res := <-ci:
		value := []*Report.Command_Param(res)[0].Npara
		var status uint64
		if value != NOT_ONLINE {
			value = []*Report.Command_Param(res)[2].Npara
			status = []*Report.Command_Param(res)[1].Npara
		}
		jsonRes, _ = json.Marshal(map[string]uint64{"result": value, "status": status})

	case <-time.After(time.Duration(delay) * time.Second):
		log.Println("res : 超时")
		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))
	GetHttpRouter().DelRequest(chanKey)

}
