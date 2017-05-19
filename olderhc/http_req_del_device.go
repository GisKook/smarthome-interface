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

//删除设备
/*
CMT_REQ_DEL_DEVICE
  参数 1.DeviceID(uint64 ieee)
    CMT_REP_DEL_DEVICE
	  参数 1.result(uint8 0.success 1.fail)
*/
func DelDeviceHandler(w http.ResponseWriter, r *http.Request) {

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
			log.Println("sorry, server break down!")
		}
	}()

	sGatewayid := r.Form["gatewayid"][0]
	tid := base.Macaddr2uint64(sGatewayid)
	sDeviceid := r.Form["deviceid"][0]
	deviceid := base.Deviceid2uint64(sDeviceid)

	// 请求参数不正确
	if sGatewayid == "" || len(sGatewayid) != 12 || sDeviceid == "" || len(sDeviceid) != 16 {
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
			Type: Report.Command_CMT_REQ_DEL_DEVICE,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT64,
					Npara: uint64(deviceid),
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)

	//log.Printf("control command request: ", reqdata)

	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	log.Println("http" + chanKey)

	ci := GetHttpRouter().SendRequest(chanKey)
	Send(reqdata)

	select {
	case res := <-ci:
		value := []*Report.Command_Param(res)[0].Npara
		if value != NOT_ONLINE {
			value = []*Report.Command_Param(res)[1].Npara
		}
		jsonRes, _ = json.Marshal(map[string]uint64{"result": value})

	case <-time.After(time.Duration(delay) * time.Second):
		log.Println("res : 超时")
		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	GetHttpRouter().DelRequest(chanKey)
	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))

}
