// 消息解析模块
package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

// 解析器维护的数据结构
// FIXME 程序会自动校正异常数据
type MessageParser struct {
	MessageLen    int    // 用来存储数据长度所占的空间：1,2,4字节，defalut=2
	MaxMessageLen uint32 // 数据最大长度
	MinMessageLen uint32 // 数据最小长度
	LittleEndian  bool   // 大小端（网络字节序都是大端模式，x86架构的主机都是小端模式）
}

// 构建一个消息解析器
// FIXME 本版本不开放此配置
func NewMessageParser() *MessageParser {
	fmt.Println("构建一个消息解析器")
	newMsg := new(MessageParser)
	newMsg.MessageLen = 2
	newMsg.MinMessageLen = 1
	newMsg.MaxMessageLen = 4
	newMsg.LittleEndian = false
	return newMsg
}

// 设置参数
func (msgParser *MessageParser) SetMsgLen(MessageLen int, MaxMessageLen uint32, MinMessageLen uint32) {
	if MessageLen == 1 || MessageLen == 2 || MessageLen == 4 {
		msgParser.MessageLen = MessageLen
	} else {
		msgParser.MessageLen = 2
	}
	if MinMessageLen != 0 {
		msgParser.MinMessageLen = MinMessageLen
	}
	if MaxMessageLen != 0 {
		msgParser.MaxMessageLen = MaxMessageLen
	}
	var max uint32
	switch msgParser.MessageLen {
	case 1:
		max = math.MaxUint8
	case 2:
		max = math.MaxUint16
	case 4:
		max = math.MaxUint32
	}
	if msgParser.MinMessageLen > max {
		msgParser.MinMessageLen = max
	}
	if msgParser.MaxMessageLen > max {
		msgParser.MaxMessageLen = max
	}
}

// 读取消息
func (msgParser *MessageParser) Read(conn *Conn) ([]byte, error) {
	var b [4]byte
	bufMsgLen := b[:msgParser.MessageLen]
	// 先获取数据长度
	if _, err := io.ReadFull(conn, bufMsgLen); err != nil {
		return nil, err
	}
	// 解析长度数据
	var msgLen uint32
	switch msgParser.MessageLen {
	// 单字节不需要处理大小端模式
	case 1:
		msgLen = uint32(bufMsgLen[0])
	// 多字节需要处理大小端模式
	case 2:
		if msgParser.LittleEndian {
			msgLen = uint32(binary.LittleEndian.Uint16(bufMsgLen))
		} else {
			msgLen = uint32(binary.BigEndian.Uint16(bufMsgLen))
		}
	case 4:
		if msgParser.LittleEndian {
			msgLen = binary.LittleEndian.Uint32(bufMsgLen)
		} else {
			msgLen = binary.BigEndian.Uint32(bufMsgLen)
		}
	}
	// 检查长度
	switch {
	case msgLen > msgParser.MaxMessageLen:
		return nil, errors.New("message too long")
	case msgLen < msgParser.MinMessageLen:
		return nil, errors.New("message too short")
	}
	// 这里才是真正获取消息体
	msgData := make([]byte, msgLen)
	if _, err := io.ReadFull(conn, msgData); err != nil {
		return nil, err
	}
	return msgData, nil
}

// 写数据
func (msgParser *MessageParser) Write(conn *Conn, args ...[]byte) error {
	// 获取数据长度
	var msgLen uint32
	for _, value := range args {
		msgLen += uint32(len(value))
	}
	// 检查长度
	switch {
	case msgLen > msgParser.MaxMessageLen:
		return errors.New("message too long")
	case msgLen < msgParser.MinMessageLen:
		return errors.New("message too short")
	}
	msg := make([]byte, uint32(msgParser.MessageLen)+msgLen)
	// 先写入消息体的长度数据
	switch msgParser.MessageLen {
	case 1:
		msg[0] = byte(msgLen)
	case 2:
		if msgParser.LittleEndian {
			binary.LittleEndian.PutUint16(msg, uint16(msgLen))
		} else {
			binary.BigEndian.PutUint16(msg, uint16(msgLen))
		}
	case 4:
		if msgParser.LittleEndian {
			binary.LittleEndian.PutUint32(msg, msgLen)
		} else {
			binary.BigEndian.PutUint32(msg, msgLen)
		}
	}
	// 最后写入编码后的数据
	length := msgParser.MessageLen
	for i := 0; i < len(args); i++ {
		copy(msg[length:], args[i])
		length += len(args[i])
	}
	conn.Write(msg)
	return nil
}
