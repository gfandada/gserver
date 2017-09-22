package gateway

import (
	"errors"
	"fmt"
	"io"

	"github.com/gfandada/gserver/network"
	Services "github.com/gfandada/gserver/services"
	Discovery "github.com/gfandada/gserver/services/discovery"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	ERROR_STREAM_NOT_OPEN = errors.New("stream not opened yet")
)

type router struct {
	sess   *Session
	config *network.Config
}

func (r *router) run(sync chan struct{}) {
	conn := Discovery.GetService("game")
	cli := network.NewServiceClient(conn)
	ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userid": fmt.Sprint(r.sess.UserId)}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		return
	}
	r.sess.Stream = stream
	sync <- struct{}{}
	for {
		in, err := r.sess.Stream.Recv()
		if err == io.EOF { // 流正常关闭
			return
		}
		if err != nil { // 出错
			return
		}
		select {
		case r.sess.MQ <- *in:
		case <-r.sess.Die:
			return
		}
	}
}

// gateway local router:[0,1999]
// gateway remote router:[2000,+}
func (r *router) router(id uint16, msg []byte) []byte {
	if id >= 2000 {
		if err := r.remoteroute(msg); err != nil {
			r.sess.Flag |= SESS_AUTHFAILED
			return Services.NewInError(err)
		}
		return nil
	}
	return r.localrouter(id, msg)
}

func (r *router) localrouter(id uint16, msg []byte) []byte {
	return r.handler(id, msg)
}

func (r *router) remoteroute(p []byte) error {
	frame := &network.Data_Frame{
		Type:    network.Data_Message,
		Message: p,
	}
	if r.sess.Stream == nil {
		return ERROR_STREAM_NOT_OPEN
	}
	if err := r.sess.Stream.Send(frame); err != nil {
		return err
	}
	return nil
}

func (r *router) handler(id uint16, msg []byte) []byte {
	if hand := Services.GetHandler(id); hand != nil {
		data, err := r.config.MsgParser.Deserialize(msg)
		if err != nil {
			return Services.NewInError(err)
		}
		return r.ackhandler(hand([]interface{}{data}))
	} else {
		r.sess.Flag |= SESS_AUTHFAILED
	}
	return nil
}

func (r *router) ackhandler(ack []interface{}) []byte {
	if ack == nil {
		return nil
	}
	data, err := r.config.MsgParser.Serialize(ack[0].(network.RawMessage))
	if err != nil {
		return Services.NewInError(err)
	}
	switch len(ack) {
	case 1:
		return data
	case 2:
		// update session
		r.sess.UserId = ack[0].(int32)
		return data
	default:
	}
	return nil
}

// 构建gateway->game路由器
// @params sess:会话
// @return GateRecv处理器
func startRouter(sess *Session, config *network.Config) *router {
	if sess == nil {
		return nil
	}
	r := &router{
		sess:   sess,
		config: config,
	}
	sync := make(chan struct{}, 1)
	go r.run(sync)
	<-sync
	close(sync)
	return r
}
