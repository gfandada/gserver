// 网关接收器
package gateway

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
	Services "github.com/gfandada/gserver/services"
)

type gaterecv struct {
	sess          *Session
	in            <-chan []byte
	out           *gatesend
	config        *network.Config
	one_min_timer <-chan time.Time
	router        *router
}

func (gr *gaterecv) run() {
	sess := gr.sess
	out := gr.out
	one_min_timer := time.After(time.Minute)
	defer func() {
		close(sess.Die)
		if sess.Stream != nil {
			sess.Stream.CloseSend()
		}
	}()
	for {
		select {
		case data, ok := <-gr.in:
			if !ok {
				return
			}
			gr.clientToGate(data)
		case frame := <-sess.MQ:
			switch frame.Type {
			case network.Data_Message:
				out.send(frame.Message)
			case network.Data_Kick:
				sess.Flag |= SESS_AUTHFAILED
			}
		case <-one_min_timer:
			gr.one_timer_work()
		case <-sess.Die:
			sess.Flag |= SESS_AUTHFAILED
		}
		if sess.Flag&SESS_AUTHFAILED != 0 {
			return
		}
	}
}

func (gr *gaterecv) clientToGate(data []byte) {
	defer func() {
		if r := recover(); r != nil {
			switch reflect.TypeOf(r).Name() {
			case "int":
				gr.out.send(Services.NewLogicError(r.(int)))
			default:
				gr.out.send(Services.NewInError(fmt.Errorf("%v", r)))
			}
		}
	}()
	gr.sess.PacketCount++
	gr.sess.PacketCountOneMin++
	gr.sess.PacketTime = time.Now()
	gr.sess.LastPacketTime = gr.sess.PacketTime
	id, msg, err := gr.check(data)
	if err == nil {
		if result := gr.router.router(id, msg); result != nil {
			gr.out.send(result)
		}
	} else {
		gr.out.send(Services.NewInError(err))
	}
	return
}

func (gr *gaterecv) one_timer_work() {
	defer func() {
		gr.sess.PacketCountOneMin = 0
		gr.one_min_timer = time.After(time.Minute)
	}()
	if gr.sess.PacketCountOneMin > gr.config.Rpm {
		gr.sess.Flag |= SESS_AUTHFAILED
	}
}

func (gr *gaterecv) check(data []byte) (uint16, []byte, error) {
	seq, id, msg, err := gr.config.Parser.ReadBodyFull(data)
	if err != nil {
		gr.sess.Flag |= SESS_AUTHFAILED
		return 0, nil, errors.New("ReadBodyFull is wrong")
	}
	if seq != gr.sess.PacketCount {
		gr.sess.Flag |= SESS_AUTHFAILED
		return 0, nil, errors.New("seq is wrong")
	}
	return id, msg, nil
}

// 构建GateRecv处理器
// @params sess:会话  in:处理client->gateway out:处理gateway->client config:配置参数
// @return GateRecv处理器
func startRecver(sess *Session, in <-chan []byte, out *gatesend, config *network.Config) *gaterecv {
	if sess == nil || out == nil || config == nil {
		return nil
	}
	sess.MQ = make(chan network.Data_Frame, config.AsyncMQ)
	sess.ConnectTime = time.Now()
	sess.LastPacketTime = time.Now()
	one_min_timer := time.After(time.Minute)
	gr := &gaterecv{
		sess:          sess,
		in:            in,
		out:           out,
		config:        config,
		one_min_timer: one_min_timer,
		router:        startRouter(sess, config),
	}
	go gr.run()
	logger.Debug(fmt.Sprintf("recver run %d", sess.UserId))
	return gr
}
