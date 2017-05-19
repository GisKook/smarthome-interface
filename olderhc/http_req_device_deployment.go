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

//布防撤销指令处理
/*
CMT_REQ_DEPLOYMENT
  B.参数 1.deviceid(uint64) 2. endpoint(uint8) 3.armmodel(uint8 0 arm 1 disarm 2 armtime)
           4.armstarttime_hour(uint8 0-24) 5.armstarttime_min(uint8 0-60)
		            6.armendtime_hour(uint8 0-24) 7.armendtime_min(uint8 0-60)
*/
func DefenceCancleHandler(w http.ResponseWriter, r *http.Request) {
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
	nArmmodel, _ := strconv.Atoi(r.Form["armmodel"][0])
	//nArm, _ := strconv.Atoi(r.Form["arm"][0])
	nArmstarttime_hour, _ := strconv.Atoi(r.Form["armstarttime_hour"][0])
	nArmstarttime_min, _ := strconv.Atoi(r.Form["armstarttime_min"][0])
	nArmendtime_hour, _ := strconv.Atoi(r.Form["armendtime_hour"][0])
	nArmendtime_min, _ := strconv.Atoi(r.Form["armendtime_min"][0])

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
			Type: Report.Command_CMT_REQ_DEPLOYMENT,
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
					Npara: uint64(nArmmodel),
				},
				//				&Report.Command_Param{
				//					Type:  Report.Command_Param_UINT8,
				//					Npara: uint64(nArm),
				//				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nArmstarttime_hour),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nArmstarttime_min),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nArmendtime_hour),
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nArmendtime_min),
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)

	log.Println("send nsq mac " + strconv.Itoa(int(tid)))
	log.Println("send nsq ser " + strconv.Itoa(int(s)))
	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	log.Println("send nsq " + chanKey)
	ci := GetHttpRouter().SendRequest(chanKey)
	Send(reqdata)

	select {
	case res := <-ci:
		log.Println(res)
		value := []*Report.Command_Param(res)[0].Npara
		//log.Println(value)
		if value != NOT_ONLINE {
			value = []*Report.Command_Param(res)[6].Npara
		}
		jsonRes, _ = json.Marshal(map[string]uint64{"result": value})
	case <-time.After(time.Duration(delay) * time.Second):
		log.Println("res : 超时")
		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))
	GetHttpRouter().DelRequest(chanKey)
}
