package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 处理器的数据结构
type MsgManager struct {
	MsgMap       map[uint16]*MessageInfo // id池：主要用于识别id对应的pb结构
	LittleEndian bool
}

// 消息
type MessageInfo struct {
	MsgType reflect.Type
}

// 构建一个新的消息处理器
func NewMsgManager() *MsgManager {
	manager := new(MsgManager)
	manager.MsgMap = make(map[uint16]*MessageInfo)
	manager.LittleEndian = false
	return manager
}

/******************************实现了imessage接口*****************************/

func (msgManager *MsgManager) NewIMessage() Imessage {
	return &MsgManager{MsgMap: msgManager.MsgMap, LittleEndian: msgManager.LittleEndian}
}

func (msgManager *MsgManager) Register(rawM *RawMessage) error {
	if _, ok := msgManager.MsgMap[rawM.MsgId]; ok {
		return fmt.Errorf("msg has registered", rawM.MsgId)
	}
	if len(msgManager.MsgMap) >= math.MaxUint16 {
		return fmt.Errorf("too many protobuf messages (max = %v)", math.MaxUint16)
	}
	newMessage := new(MessageInfo)
	newMessage.MsgType = reflect.TypeOf(rawM.MsgData.(proto.Message))
	msgManager.MsgMap[rawM.MsgId] = newMessage
	return nil
}

func (msgManager *MsgManager) UnRegister(rawM *RawMessage) {
	delete(msgManager.MsgMap, rawM.MsgId)
}

// for id+message
func (msgManager *MsgManager) Serialize(rawM RawMessage) ([]byte, error) {
	if _, ok := msgManager.MsgMap[rawM.MsgId]; !ok {
		return nil, errors.New("message has not registered")
	}
	if rawM.MsgRaw == nil {
		rawId := make([]byte, 2)
		if msgManager.LittleEndian {
			binary.LittleEndian.PutUint16(rawId, rawM.MsgId)
		} else {
			binary.BigEndian.PutUint16(rawId, rawM.MsgId)
		}
		data, err := proto.Marshal(rawM.MsgData.(proto.Message))
		if err != nil {
			return nil, err
		}
		c := make([]byte, 2+len(data))
		copy(c, rawId)
		copy(c[len(rawId):], data)
		return c, err
	} else {
		return rawM.MsgRaw, nil
	}
}

// for id+message
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
		//if info.MsgClient != nil {
		msg := reflect.New(info.MsgType.Elem()).Interface()
		err := proto.Unmarshal(data[2:], msg.(proto.Message))
		return &RawMessage{MsgId: id, MsgData: msg, MsgRaw: data}, err
		//}
		//return &RawMessage{MsgRaw: data}, nil
	}
	return &RawMessage{}, fmt.Errorf("message %d has not registered", id)
}
