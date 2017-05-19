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

func ReadDeployMentStatusHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var jsonRes []byte //返回内容
	//处理panic
	defer func() {
		if x := recover(); x != nil {
			jsonRes, _ = json.Marshal(map[string]byte{"result": SERVER_FAILED})
			log.Println(string(jsonRes))
			fmt.Fprint(w, string(jsonRes))
			log.Println("sorry, server break down!ReadDeployMentStatusHandler")
		}
	}()

	sGatewayid := r.Form["gatewayid"][0]
	tid := base.Macaddr2uint64(sGatewayid)
	sDeviceid := r.Form["deviceid"][0]
	deviceid := base.Deviceid2uint64(sDeviceid)
	log.Println("ReadDeployMentStatusHandler" + sDeviceid)
	sEndpoint := r.Form["endpoint"][0]
	nEndpoint, _ := strconv.Atoi(sEndpoint)
	log.Println("ReadDeployMentStatusHandler" + sEndpoint)
	serialnum := r.Form["serialid"][0]
	nSerialID, _ := strconv.Atoi(serialnum)
	//nSerialID := GenerateKey(tid, nSerialID)

	// 请求参数不正确
	if sGatewayid == "" || len(sGatewayid) != 12 || sDeviceid == "" || len(sDeviceid) != 16 || sEndpoint == "" {
		jsonRes, _ = json.Marshal(map[string]byte{"result": BAD_PARAMETER})
		log.Println(string(jsonRes))
		fmt.Fprint(w, string(jsonRes))
		log.Println("wrong arguments!")
		return
	} else {
		log.Println("arguents correct!")
	}

	//构造指令内容
	req := &Report.ControlReport{
		Tid:          uint64(tid),
		SerialNumber: uint32(nSerialID),
		Command: &Report.Command{
			Type: Report.Command_CMT_REQ_READ_DEPLOYMENT_STATUS,
			Paras: []*Report.Command_Param{
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT64,
					Npara: deviceid,
				},
				&Report.Command_Param{
					Type:  Report.Command_Param_UINT8,
					Npara: uint64(nEndpoint),
				},
			},
		},
	}
	reqdata, _ := proto.Marshal(req)
	log.Println("before send nsq mac ")

	log.Println("send nsq mac " + strconv.Itoa(int(tid)))
	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(nSerialID))
	log.Println("send nsq " + chanKey)
	ci := GetHttpRouter().SendRequest(chanKey)
	Send(reqdata)

	select {
	case res := <-ci:
		log.Println(res)
		value := []*Report.Command_Param(res)[0].Npara
		//log.Println(value)
		if value != NOT_ONLINE {
			armmodel := []*Report.Command_Param(res)[1].Npara
			start_hour := []*Report.Command_Param(res)[2].Npara
			start_min := []*Report.Command_Param(res)[3].Npara
			end_hour := []*Report.Command_Param(res)[4].Npara
			end_min := []*Report.Command_Param(res)[5].Npara
			value = []*Report.Command_Param(res)[6].Npara
			jsonRes, _ = json.Marshal(map[string]uint64{"result": value,
				"armmodel": armmodel, "start_hour": start_hour, "start_min": start_min,
				"end_hour": end_hour, "end_min": end_min})
		} else {
			jsonRes, _ = json.Marshal(map[string]uint64{"result": value})
		}
	case <-time.After(time.Duration(delay) * time.Second):
		log.Println("res : 超时")
		jsonRes, _ = json.Marshal(map[string]byte{"result": CONTROL_FAILURE})
	}

	log.Println(string(jsonRes))
	fmt.Fprint(w, string(jsonRes))
	GetHttpRouter().DelRequest(chanKey)
}
