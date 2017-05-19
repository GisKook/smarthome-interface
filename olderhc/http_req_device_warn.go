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

//设备报警指令处理
/*
CMT_REQ_DEVICE_WARN
  参数 1.DeviceID(uint64) 2.endpoint(uint8) 3.warningduration(uint8报警时长，s)
           4.WarningMode(uint8 取值0~6) 5.storebe(uint8 0-1) 6.sirenlevel(uint8 0-3)
		            7.strobelevel(uint8 0-3) 8.strobedutycycle(uint8 01)
					  返回结果  总是返回发送成功.
*/
func DevWarnHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	now := time.Now()
	var s int32 = int32(now.Unix()) //% 255
	var jsonRes []byte              //返回内容
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
	nWarningduration, _ := strconv.Atoi(r.Form["warningduration"][0])
	nWarningmode, _ := strconv.Atoi(r.Form["warningmode"][0])
	nStorebe, _ := strconv.Atoi(r.Form["storebe"][0])
	nSirenlevel, _ := strconv.Atoi(r.Form["sirenlevel"][0])
	nStrobelevel, _ := strconv.Atoi(r.Form["strobelevel"][0])
	nStrobedutycycle, _ := strconv.Atoi(r.Form["strobedutycycle"][0])

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
			Type: Report.Command_CMT_REQ_DEVICE_WARN,
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
					Npara: uint64(nWarningduration),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nWarningmode),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nStorebe),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nSirenlevel),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nStrobelevel),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nStrobedutycycle),
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)
	Send(reqdata)

	jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_SUCCESS})
	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))
}
