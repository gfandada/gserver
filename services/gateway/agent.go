// 通用websocket网关
package gateway

import (
	"time"

	"github.com/gfandada/gserver/network"
)

type Agent struct {
	configdata *network.Config // 配置数据
	conn       network.Iconn   // 当前会话
}

func (agent *Agent) Start(conn network.Iconn) {
	defer conn.Close()
	agent.conn = conn
	config := agent.configdata
	in := make(chan []byte)
	defer close(in)
	var sess Session
	sess.Die = make(chan struct{})
	if sender := startSender(conn, &sess, in, config); sender == nil {
		close(sess.Die)
		return
	}
	for {
		conn.SetReadDeadline(time.Now().Add(time.Duration(config.ReadDeadline) * time.Second))
		conn.SetWriteDeadline(time.Now().Add(time.Duration(config.WriteDeadline) * time.Second))
		data, err := conn.ReadMsg()
		if err != nil {
			return
		}
		payload, errp := config.Parser.ReadBody(data)
		if errp != nil {
			return
		}
		select {
		case in <- payload:
		case <-sess.Die:
			return
		}
	}
	agent.Close()
}

func (agent *Agent) NewIagent() network.Iagent {
	return &Agent{configdata: agent.configdata}
}

func (agent *Agent) Close() {

}

func (agent *Agent) GetUserData() interface{} {
	return nil
}

func (agent *Agent) SetUserData(data interface{}) {

}
