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

//查询设备属性指令处理
/*
CMT_REQ_DEVICE_ATTR
  参数 1.DeviceID(uint64 ieee) 2.endpoint(uint8)
    CMT_REP_DEVICE_ATTR
	  参数 1.shortaddr(uint16) 2.profileid(uint16) 3.zclversion (uint8) 4.applicationversion(uint8)
	      5.stackversion(uint8) 6.hwversion(uint8) 7.manufacturename(utf8 string) 8.modelidentifier(utf8 string)
		      9.datecode(utf8 string) 10.powersource(uint8)
*/

type DeviceAttr struct {
	Shortaddr          uint64
	Profileid          uint64
	Zclversionst       uint64
	Applicationversion uint64
	Stackversion       uint64
	Hwversion          uint64
	Manufacturename    string
	Modelidentifier    string
	Datecode           string
	Powersource        uint64
}

func GetDevAttributeHandler(w http.ResponseWriter, r *http.Request) {
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
	sEndpoint := r.Form["endpoint"][0]
	nEndpoint, _ := strconv.Atoi(sEndpoint)

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
			Type: Report.Command_CMT_REQ_DEVICE_ATTR,
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

	//log.Printf("control command request: ", reqdata)

	chanKey := strconv.Itoa(int(tid)) + strconv.Itoa(int(s))
	log.Println("http" + chanKey)

	ci := GetHttpRouter().SendRequest(chanKey)
	Send(reqdata)

	select {
	case res := <-ci:
		//var paras []*Report.Command_Param
		//paras = res
		//	log.Println(res)
		value := []*Report.Command_Param(res)[0].Npara
		if value != NOT_ONLINE {
			//Shortaddr          uint64
			//Profileid          uint64
			//Zclversionst       uint64
			//Applicationversion uint64
			//Stackversion       uint64
			//Hwversion          uint64
			//Manufacturename    string
			//Modelidentifier    string
			//Datecode           string
			//Powersource        uint64
			var attr DeviceAttr
			attr.Shortaddr = []*Report.Command_Param(res)[1].Npara
			attr.Profileid = []*Report.Command_Param(res)[2].Npara
			attr.Zclversionst = []*Report.Command_Param(res)[3].Npara
			attr.Applicationversion = []*Report.Command_Param(res)[4].Npara
			attr.Stackversion = []*Report.Command_Param(res)[5].Npara
			attr.Hwversion = []*Report.Command_Param(res)[6].Npara
			attr.Manufacturename = []*Report.Command_Param(res)[7].Strpara
			attr.Modelidentifier = []*Report.Command_Param(res)[8].Strpara
			attr.Datecode = []*Report.Command_Param(res)[9].Strpara
			attr.Powersource = []*Report.Command_Param(res)[10].Npara
			log.Println(attr)
			jsonRes, _ = json.Marshal(attr)
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
