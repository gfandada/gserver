// 基于pb来处理消息体
package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"lib/chanrpc"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
)

// 处理器的数据结构
type Processor struct {
	LittleEndian bool                    // FIXME 大小端模式 注意需要和消息解析器保持一致
	MsgInfo      []*MessageInfo          // FIXME 消息体容器 注意只允许是指针
	MsgId        map[reflect.Type]uint16 // FIXME 消息id容器 注意id是个unit16
}

// 本结构不用来标识消息本身的数据结构
type MessageInfo struct {
	MsgType       reflect.Type    // 消息类型
	MsgRouter     *chanrpc.Server // TODO 基于chan的消息路由器
	MsgHandler    MessageHandler  // 消息处理器（需要使用特定编码处理）
	MsgRawHandler MessageHandler  // 消息处理器（不需要使用特定编码处理）
}

// 定义的消息回调函数
type MessageHandler func([]interface{})

// 这里是定义的转码后的消息体
type MsgRaw struct {
	MsgId      uint16 // 消息id
	MsgRawData []byte // 消息的真实数据
}

// 构建一个新的消息处理器
func NewProcessor() *Processor {
	processor := new(Processor)
	// 这里和MessageParser一样保持大端模式
	processor.LittleEndian = false
	processor.MsgId = make(map[reflect.Type]uint16)
	return processor
}

// 设置字节序
// FIXME 非携程安全
func (processor *Processor) SetByteOrder(littleEndian bool) {
	processor.LittleEndian = littleEndian
}

// 注册消息
// FIXME 非携程安全
func (processor *Processor) Register(msg proto.Message) uint16 {
	// 反射出消息类型
	msgType := reflect.TypeOf(msg)
	// 检查下消息的类型
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		fmt.Println("protobuf message pointer required")
	}
	// 检查是否已经被系统注册的消息
	if _, ok := processor.MsgId[msgType]; ok {
		fmt.Printf("message %s is already registered\n", msgType)
	}
	// 检查消息长度
	if len(processor.MsgInfo) >= math.MaxUint16 {
		fmt.Printf("too many protobuf messages (max = %v)", math.MaxUint16)
	}
	// 构建消息
	message := new(MessageInfo)
	message.MsgType = msgType
	// 在消息体容器中增加这些消息
	processor.MsgInfo = append(processor.MsgInfo, message)
	// 返回从0开始的id给接口调用者
	id := uint16(len(processor.MsgInfo) - 1)
	// 在id容器中增加这些id标识
	processor.MsgId[msgType] = id
	return id
}

// 设置消息路由器
// FIXME 在routing和序列化、反序列化的过程中禁止调用此方案
func (processor *Processor) SetRouter(msg proto.Message, msgRouter *chanrpc.Server) {
	msgType := reflect.TypeOf(msg)
	// 检查消息是否被注册
	id, ok := processor.MsgId[msgType]
	if !ok {
		fmt.Printf("SetRouter message %s not registered", msgType)
		return
	}
	processor.MsgInfo[id].MsgRouter = msgRouter
}

// 设置消息处理器
// FIXME 在routing和序列化、反序列化的过程中禁止调用此方案
func (processor *Processor) SetHandler(msg proto.Message, msgHandler MessageHandler) {
	msgType := reflect.TypeOf(msg)
	// 检查消息是否被注册
	id, ok := processor.MsgId[msgType]
	if !ok {
		fmt.Printf("SetHandler message %s not registered", msgType)
		return
	}
	processor.MsgInfo[id].MsgHandler = msgHandler
}

// 设置真实消息体的处理器
// FIXME 在routing和序列化、反序列化的过程中禁止调用此方案
func (processor *Processor) SetRawHandler(id uint16, msgRawHandler MessageHandler) {
	if id >= uint16(len(processor.MsgInfo)) {
		fmt.Printf("message id %v not registered", id)
	}
	processor.MsgInfo[id].MsgRawHandler = msgRawHandler
}

// 批量处理消息类型，使用指定的回调方法
func (processor *Processor) Range(fun func(id uint16, t reflect.Type)) {
	for index, info := range processor.MsgInfo {
		fun(uint16(index), info.MsgType)
	}
}

/******************************实现了imessage接口*****************************/

// 序列化
func (processor *Processor) Serialize(msg interface{}) ([][]byte, error) {
	// 反射消息的类型
	msgType := reflect.TypeOf(msg)
	// 检查消息是否被注册了
	id, ok := processor.MsgId[msgType]
	if !ok {
		err := fmt.Errorf("message %s not registered", msgType)
		fmt.Printf("message %s not registered\n", msgType)
		return nil, err
	}
	// 序列化消息id
	rawId := make([]byte, 2)
	if processor.LittleEndian {
		binary.LittleEndian.PutUint16(rawId, id)
	} else {
		binary.BigEndian.PutUint16(rawId, id)
	}
	// 序列化消息体
	data, err := proto.Marshal(msg.(proto.Message))
	return [][]byte{rawId, data}, err
}

// 反序列化接口
func (processor *Processor) Deserialize(data []byte) (interface{}, error) {
	if len(data) < 2 {
		return nil, errors.New("protobuf data too short")
	}
	// 反序列消息id
	var id uint16
	if processor.LittleEndian {
		id = binary.LittleEndian.Uint16(data)
	} else {
		id = binary.BigEndian.Uint16(data)
	}
	if id >= uint16(len(processor.MsgInfo)) {
		err := fmt.Errorf("message id %v not registered", id)
		fmt.Printf("message id %v not registered\n", id)
		return nil, err
	}
	// 反序列化消息体
	// TODO
	info := processor.MsgInfo[id]
	if info.MsgRawHandler != nil {
		return MsgRaw{id, data[2:]}, nil
	} else {
		msg := reflect.New(info.MsgType.Elem()).Interface()
		return msg, proto.UnmarshalMerge(data[2:], msg.(proto.Message))
	}
}

// 路由器
func (processor *Processor) Route(msg interface{}, userData interface{}) error {
	// 必要的检查
	if msgRaw, ok := msg.(MsgRaw); ok {
		if msgRaw.MsgId >= uint16(len(processor.MsgInfo)) {
			err := fmt.Errorf("message msg %v not registered", msg)
			fmt.Printf("Route message msg %v not registered\n", msg)
			return err
		}
		info := processor.MsgInfo[msgRaw.MsgId]
		if info.MsgRawHandler != nil {
			info.MsgRawHandler([]interface{}{msgRaw.MsgId, msgRaw.MsgRawData, userData})
		}
		return nil
	}
	// protobuf
	msgType := reflect.TypeOf(msg)
	id, ok := processor.MsgId[msgType]
	if !ok {
		return fmt.Errorf("message %s not registered", msgType)
	}
	info := processor.MsgInfo[id]
	if info.MsgHandler != nil {
		info.MsgHandler([]interface{}{msg, userData})
	}
	// 如果有路由器，则
	if info.MsgRouter != nil {
		info.MsgRouter.Go(msgType, msg, userData)
	}
	return nil
}
