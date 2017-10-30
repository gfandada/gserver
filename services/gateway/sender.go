// 网关发送器：gateway->client
package gateway

import (
	"encoding/binary"
	"time"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
)

const (
	RECV_NORMAL = "recvnormal" // from client->gateway
	RECV_CTRL   = "recvctrl"   // from client->gateway
	SEND_NORMAL = "sendnormal" // from gateway->client
	SEND_CTRL   = "sendctrl"   // from gateway->client
)

type gatesend struct {
	die     <-chan struct{}
	pending chan []byte
	conn    network.Iconn
	recver  *gaterecv
	config  *network.Config
}

func (gtc *gatesend) send(data []byte) {
	if data == nil {
		return
	}
	select {
	case gtc.pending <- data:
	default: // default drop data
	}
}

func (gtc *gatesend) run() {
	for {
		select {
		case data := <-gtc.pending:
			gtc.raw_send(data)
		case <-gtc.die:
			return
		}
	}
}

func (gtc *gatesend) raw_send(data []byte) {
	msg, err := gtc.config.Parser.Write(data)
	// FIXME for debug
	id := binary.BigEndian.Uint16(data)
	if err == nil {
		// FIXME ignore failed
		gtc.conn.SetWriteDeadline(time.Now().Add(time.Duration(gtc.config.WriteDeadline) * time.Second))
		err = gtc.conn.WriteMsg(msg)
		logger.Debug("send msgid %d to %v datalen %d err %v", id, gtc.conn.RemoteAddr(), len(msg), err)
	}
}

// 构建GateToClient处理器
// @params conn:会话  die:控制器 pendingnum:排队上限
// @return GateToClient处理器
func startSender(conn network.Iconn, sess *Session, in <-chan []byte, config *network.Config) *gatesend {
	if conn == nil {
		return nil
	}
	cgs := &gatesend{
		conn:    conn,
		pending: make(chan []byte, config.PendingNum),
		die:     sess.Die,
		config:  config,
	}
	go cgs.run()
	cgs.recver = startRecver(sess, in, cgs, config)
	return cgs
}
