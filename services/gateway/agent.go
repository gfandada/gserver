// 通用的代理器
package gateway

import (
	"fmt"
	"time"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
)

type Agent struct {
	configdata *network.Config // 配置数据
	conn       network.Iconn   // 当前会话
}

func (agent *Agent) Start(conn network.Iconn) {
	defer conn.Close()
	// for goroutine safe
	agent.configdata.Parser = agent.configdata.Parser.NewMessageParser()
	agent.conn = conn
	config := agent.configdata
	in := make(chan []byte)
	defer close(in)
	defer agent.Close()
	var sess Session
	sess.Die = make(chan struct{})
	logger.Debug(fmt.Sprintf("agent run %v", agent.conn.RemoteAddr()))
	if sender := startSender(conn, &sess, in, config); sender == nil {
		close(sess.Die)
		logger.Error(fmt.Sprintf("agent run sender nil"))
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
}

func (agent *Agent) NewIagent() network.Iagent {
	return &Agent{configdata: agent.configdata}
}

func (agent *Agent) Close() {
	logger.Debug(fmt.Sprintf("agent close %v", agent.conn.RemoteAddr()))
}

func (agent *Agent) GetUserData() interface{} {
	return nil
}

func (agent *Agent) SetUserData(data interface{}) {

}
