package olderhc

import (
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/giskook/smarthome-interface/olderhc/pbgo"
	"log"
)

// 是否在线返回结果
const (
	NOT_ONLINE = 252
	YES_ONLINE = 251
)

// 控制返回结果
const (
	CONTROL_SUCCESS = 0
	CONTROL_FAILURE = 255
	BAD_PARAMETER   = 253
	SERVER_FAILED   = 254
)

type basic struct {
	commandType int32
	mac         string //需确定das处，mac存储类型
	serialNum   int32
}

// 指令
type command struct {
	basic
	topic string
}

// 指令返回结果
type result struct {
	basic
	info  byte
	paras []*Report.Command_Param
}

// 将指令内容转换为字节数组
func (c *command) GetBytes() []byte {
	bytes := make([]byte, 0, 18)
	bytes = append(bytes, Int32ToBytes(c.basic.commandType)...)
	//bytes = append(bytes, c.basic.commandType)
	bytes = append(bytes, []byte(c.basic.mac)...)
	//bytes = append(bytes, Int32ToBytes(c.basic.mac)...)
	bytes = append(bytes, Int32ToBytes(c.basic.serialNum)...)
	bytes = append(bytes, []byte(c.topic)...)
	return bytes
}

// 将接收到的字节数组转换为结果对象
func BuildResult(message []byte) result {

	data := message
	report := &Report.ControlReport{}
	err := proto.Unmarshal(data, report)
	if err != nil {
		log.Println("unmarshal error")
	}

	mac := report.Tid
	serialNum := int32(report.SerialNumber)
	pbcommandType := report.GetCommand().Type
	var info uint8
	var paras []*Report.Command_Param
	paras = report.GetCommand().GetParas()
	/*
		    switch pbcommandType {
			case Report.Command_CMT_REPLOGIN:
		        paras =report.GetCommand().GetParas()
		        paraType:= paras[0].Type
		        if(paraType != Report.Command_Param_UINT8){
		            log.Println("paraType error")
		        }
		        info = uint8(paras[0].Npara)

			}*/
	// topic := fmt.Sprintf("%d", gatewayid)
	// msg := fmt.Sprintf("%d",serialnum)

	//info := message[17]
	sMac := fmt.Sprintf("%d", mac)
	res := result{basic{int32(pbcommandType), sMac, serialNum}, info, paras}
	return res
}

func Int32ToBytes(i int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}
