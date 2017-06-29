package protobuff

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lib/gservices"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 处理器的数据结构
type MsgManager struct {
	LittleEndian bool
	MsgMap       map[uint16]*MessageInfo // id池：主要用于识别id对应的pb结构
}

// 消息
type MessageInfo struct {
	MsgType    reflect.Type
	MsgHandler gservices.MessageHandler3 // 消息处理器
	MsgClient  *gservices.LocalClient    // 消息服务器
}

// 消息体
type RawMessage struct {
	MsgId   uint16
	MsgData interface{}
}

// 构建一个新的消息处理器
func NewMsgManager() *MsgManager {
	manager := new(MsgManager)
	manager.MsgMap = make(map[uint16]*MessageInfo)
	manager.LittleEndian = false
	return manager
}

// 注册消息
// FIXME 非携程安全
func (msgManager *MsgManager) RegisterMessage(rawM RawMessage, handler gservices.MessageHandler3, msgServer *gservices.LocalServer) {
	if _, ok := msgManager.MsgMap[rawM.MsgId]; ok {
		fmt.Println("msg has registered", rawM.MsgId)
		return
	}
	if len(msgManager.MsgMap) >= math.MaxUint16 {
		fmt.Println("too many protobuf messages (max = %v)", math.MaxUint16)
		return
	}
	newMessage := new(MessageInfo)
	newMessage.MsgType = reflect.TypeOf(rawM.MsgData.(proto.Message))
	newMessage.MsgHandler = handler
	newMessage.MsgClient = msgServer.NewLocalClient()
	msgManager.MsgMap[rawM.MsgId] = newMessage
}

/******************************实现了imessage接口*****************************/

func (msgManager *MsgManager) Serialize(rawM RawMessage) ([][]byte, error) {
	if _, ok := msgManager.MsgMap[rawM.MsgId]; !ok {
		return nil, errors.New("message has not registered")
	}
	rawId := make([]byte, 2)
	if msgManager.LittleEndian {
		binary.LittleEndian.PutUint16(rawId, rawM.MsgId)
	} else {
		binary.BigEndian.PutUint16(rawId, rawM.MsgId)
	}
	data, err := proto.Marshal(rawM.MsgData.(proto.Message))
	return [][]byte{rawId, data}, err
}

func (msgManager *MsgManager) Deserialize(data []byte) (*RawMessage, error) {
	if len(data) < 2 {
		return &RawMessage{}, errors.New("protobuf data too short")
	}
	var id uint16
	if msgManager.LittleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}
	if info, ok := msgManager.MsgMap[id]; ok {
		msg := reflect.New(info.MsgType.Elem()).Interface()
		err := proto.Unmarshal(data[2:], msg.(proto.Message))
		return &RawMessage{MsgId: id, MsgData: msg}, err
	}
	return &RawMessage{}, errors.New("message has not registered")
}

func (msgManager *MsgManager) Router(msg *RawMessage, userData interface{}) error {
	if info, ok := msgManager.MsgMap[msg.MsgId]; ok {
		if info.MsgClient != nil {
			info.MsgClient.Cast(&gservices.InputMessage{
				Msg:        msg.MsgId,
				F:          info.MsgHandler,
				CB:         userData.(gservices.Iack),
				Args:       []interface{}{msg.MsgData, userData},
				OutputChan: nil,
			})
			return nil
		}
	}
	return errors.New("message has not registered")
}
