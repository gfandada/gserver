// 报文解析模块
package network

/*
	client->gateway
	----------------------------
	| len | seq | id | message |
	----------------------------
	len:seq + id + message的长度，占用2个字节(uint16)
	seq:从1自增的序列号，占用4个字节(uint32)
	id:协议号，占用2个字节(uint16)
	message:业务数据，占用len-4-2字节，可以使用任意编码：pb/json等

	gateway->client
	----------------------
	| len | id | message |
	----------------------
	len:message的长度，占用2个字节(uint16)
	id:协议号，占用两个字节(uint16)
	message:业务数据，占用len-2字节，可以使用任意编码：pb/json等
*/

import (
	"encoding/binary"
	"errors"

	"github.com/gfandada/gserver/misc"
)

type MessageParser struct {
	maxMessageLen uint16 // for message
	minMessageLen uint16 // for message
	buff          []byte // for send
}

type RawMessage struct {
	MsgId   uint16
	MsgData interface{}
}

func NewMessageParser() (newMsg *MessageParser) {
	newMsg = new(MessageParser)
	newMsg.minMessageLen = 0
	newMsg.maxMessageLen = 512
	newMsg.buff = make([]byte, newMsg.maxMessageLen+4)
	return
}

func (msgParser *MessageParser) SetMsgLen(MaxMessageLen uint16, MinMessageLen uint16) {
	if MinMessageLen >= 0 {
		msgParser.minMessageLen = MinMessageLen
	}
	if MaxMessageLen >= 0 {
		msgParser.maxMessageLen = MaxMessageLen
	}
}

func (msgParser *MessageParser) NewMessageParser() *MessageParser {
	return &MessageParser{
		maxMessageLen: msgParser.maxMessageLen,
		minMessageLen: msgParser.minMessageLen,
		buff:          make([]byte, 8+msgParser.maxMessageLen),
	}
}

// 获取body(除header)
func (msgParser *MessageParser) ReadBody(data []byte) ([]byte, error) {
	size := binary.BigEndian.Uint16(data[:2])
	switch {
	case uint16(size) > msgParser.maxMessageLen+6:
		return nil, errors.New("message too long")
	case uint16(size) < msgParser.minMessageLen+6:
		return nil, errors.New("message too short")
	}
	return data[2:], nil
}

//
// 拆分body数据
// @return 数据1(序列号)，数据2(协议号)，数据3(id+message)，错误描述
func (msgParser *MessageParser) ReadBodyFull(data []byte) (uint32, uint16, []byte, error) {
	reader := misc.Reader(data)
	seq_id, err1 := reader.ReadU32()
	if err1 != nil {
		return 0, 0, nil, errors.New("read seq error")
	}
	id, err2 := reader.ReadU16()
	if err2 != nil {
		return 0, 0, nil, errors.New("read messageid error")
	}
	return seq_id, id, data[4:], nil
}

// @params data:id+message
// @return header+data
func (msgParser *MessageParser) Write(data []byte) ([]byte, error) {
	if data == nil {
		return nil, errors.New("data is too short")
	}
	size := uint16(len(data))
	buff := msgParser.buff
	if size-2 >= msgParser.minMessageLen && size-2 <= msgParser.maxMessageLen {
		binary.BigEndian.PutUint16(buff, uint16(size))
		copy(buff[2:], data)
		return buff[:2+size], nil
	}
	return nil, errors.New("data is too long")
}
