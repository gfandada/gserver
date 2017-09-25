package network

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 处理器的数据结构
type MsgManager struct {
	MsgMap map[uint16]*MessageInfo // id池：主要用于识别id对应的pb结构
}

// 消息
type MessageInfo struct {
	MsgType reflect.Type
}

// 构建一个新的消息处理器
func NewMsgManager() *MsgManager {
	manager := new(MsgManager)
	manager.MsgMap = make(map[uint16]*MessageInfo)
	return manager
}

/******************************实现了imessage接口*****************************/

// not goroutine safe
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

// not goroutine safe
func (msgManager *MsgManager) UnRegister(rawM *RawMessage) {
	delete(msgManager.MsgMap, rawM.MsgId)
}

// for id+message
// goroutine safe
func (msgManager *MsgManager) Serialize(rawM RawMessage) ([]byte, error) {
	if _, ok := msgManager.MsgMap[rawM.MsgId]; !ok {
		return nil, errors.New("message has not registered")
	}
	rawId := make([]byte, 2)
	binary.BigEndian.PutUint16(rawId, rawM.MsgId)
	data, err := proto.Marshal(rawM.MsgData.(proto.Message))
	if err != nil {
		return nil, err
	}
	return msgManager.bytesCombine(rawId, data), err
}

func (msgManager *MsgManager) bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

// for id+message
// goroutine safe
func (msgManager *MsgManager) Deserialize(data []byte) (*RawMessage, error) {
	if len(data) < 2 {
		return &RawMessage{}, errors.New("protobuf data too short")
	}
	var id uint16
	id = binary.BigEndian.Uint16(data)
	if info, ok := msgManager.MsgMap[id]; ok {
		msg := reflect.New(info.MsgType.Elem()).Interface()
		err := proto.Unmarshal(data[2:], msg.(proto.Message))
		return &RawMessage{MsgId: id, MsgData: msg}, err
	}
	return &RawMessage{}, fmt.Errorf("message %d has not registered", id)
}
