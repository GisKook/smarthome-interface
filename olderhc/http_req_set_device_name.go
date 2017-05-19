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

//设置名称指令处理
/*
CMT_REQ_SETNAME
  参数 1.DeviceID(uint64 网关使用mac device使用ieee) 2.名称(string utf8)
    CMT_REP_SETNAME
	  参数 1.result(uint8 0.success 1.fail)
*/
func SetNameHandler(w http.ResponseWriter, r *http.Request) {
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
	log.Println("deviceid int:", deviceid, "string:", sDeviceid)
	name := r.Form["name"][0]
	// 请求参数不正确
	if sGatewayid == "" || len(sGatewayid) != 12 || sDeviceid == "" || len(sDeviceid) != 16 || name == "" {
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
			Type: Report.Command_CMT_REQ_SETNAME,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT64,
					Npara: deviceid,
				},
				&Report.Command_Param{
					Type:    Report.Command_Param_STRING,
					Strpara: name,
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)

	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	log.Println("http req:" + chanKey)
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

	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))
	GetHttpRouter().DelRequest(chanKey)
}
