package olderhc

import (
	"encoding/binary"
	"github.com/golang/protobuf/proto"
    "fmt"
    "log"
    "github.com/huoyan108/dataRelayGo/olderhc/pbgo"
)

// 指令类型常量
const (
	ONLINE byte = 0

	LEFTUP_START byte = 10
	LEFTUP_STOP  byte = 11

	LEFTDOWN_START byte = 20
	LEFTDOWN_STOP  byte = 21

	OVERALLUP_START byte = 30
	OVERALLUP_STOP  byte = 31

	OVERALLDOWN_START byte = 40
	OVERALLDOWN_STOP  byte = 41

	BACKUP_START byte = 50
	BACKUP_STOP  byte = 51

	BACKDOWN_START byte = 60
	BACKDOWN_STOP  byte = 61

	LEGUP_START byte = 70
	LEGUP_STOP  byte = 71

	LEGDOWN_START byte = 80
	LEGDOWN_STOP  byte = 81
)

// 是否在线返回结果
const (
	NOT_ONLINE = iota
	YES_ONLINE
)

// 控制返回结果
const (
	CONTROL_FAILURE = 6
	CONTROL_SUCCESS
	BAD_PARAMETER = 4
	SERVER_FAILED = 5
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
	info byte
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
    paras =report.GetCommand().GetParas()
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
    sMac := fmt.Sprintf("%d",mac)
	res := result{basic{int32(pbcommandType), sMac, serialNum}, info,paras }
	return res
}

func Int32ToBytes(i int32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}
