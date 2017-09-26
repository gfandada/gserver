package service

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
	Services "github.com/gfandada/gserver/services"
	"google.golang.org/grpc/metadata"
)

const (
	DEFAULT_CH_SIZE = 16 // 默认玩家异步消息队列大小
)

var (
	ERROR_INCORRECT_FRAME_TYPE = errors.New("incorrect frame")
	ERROR_SERVICE_NOT_BIND     = errors.New("service not bind")
)

type Agent struct {
	msgParser network.Imessage
	sess      *Session
}

func (s *Agent) Stream(stream network.Service_StreamServer) error {
	sess := New()
	sess.Agent = s
	in := startRecver(stream, sess.Die)
	defer func() {
		Remove(sess.UserId)
		close(sess.Die)
	}()
	userid, err := s.getUserId(stream)
	if err != nil {
		return err
	}
	sess.UserId = userid
	s.sess = sess
	Add(userid, sess)
	logger.Debug(fmt.Sprintf("Stream run %d", s.sess.UserId))
	for {
		select {
		case frame, ok := <-in:
			if !ok {
				logger.Error(fmt.Sprintf("Stream error %v", err))
				return nil
			}
			s.handler(stream, frame)
		case frame := <-sess.MQ:
			if err := stream.Send(&frame); err != nil {
				logger.Error(fmt.Sprintf("Stream send error %v", err))
				return err
			}
		}
	}
}

func (s *Agent) getUserId(stream network.Service_StreamServer) (int32, error) {
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		return 0, ERROR_INCORRECT_FRAME_TYPE
	}
	if len(md["userid"]) == 0 {
		return 0, ERROR_INCORRECT_FRAME_TYPE
	}
	userid, err := strconv.Atoi(md["userid"][0])
	if err != nil {
		return 0, ERROR_INCORRECT_FRAME_TYPE
	}
	return int32(userid), nil
}

// for sync ipc
func (s *Agent) send(stream network.Service_StreamServer, frame *network.Data_Frame) error {
	return stream.Send(frame)
}

func (s *Agent) handler(stream network.Service_StreamServer, frame *network.Data_Frame) error {
	defer func() {
		if r := recover(); r != nil {
			switch reflect.TypeOf(r).Name() {
			case "int":
				s.send(stream, Services.NewSLogicError(r.(int)))
			default:
				s.send(stream, Services.NewSInError(fmt.Errorf("%v", r)))
			}
		}
	}()
	switch frame.Type {
	case network.Data_Message:
		return s.send(stream, s.dohandler(frame.Message))
	case network.Data_Ping:
		return s.send(stream, &network.Data_Frame{
			Type:    network.Data_Ping,
			Message: frame.Message,
		})
	default:
		return ERROR_INCORRECT_FRAME_TYPE
	}
}

// for async ipc
func (s *Agent) Send(msg network.RawMessage) {
	ackdata, err := s.msgParser.Serialize(msg)
	var data *network.Data_Frame
	if err != nil {
		data = Services.NewSInError(err)
	} else {
		data = &network.Data_Frame{
			Type:    network.Data_Message,
			Message: ackdata,
		}
	}
	s.sess.MQ <- *data
}

func (s *Agent) dohandler(data []byte) *network.Data_Frame {
	ret, err := s.msgParser.Deserialize(data)
	if err != nil {
		return Services.NewSInError(err)
	}
	if hand := Services.GetHandler(ret.MsgId); hand != nil {
		if err != nil {
			return nil
		}
		ack := s.ackhandler(hand([]interface{}{ret.MsgData}))
		if ack != nil {
			ackdata, erra := s.msgParser.Serialize(ack.(network.RawMessage))
			if erra != nil {
				return Services.NewSInError(err)
			}
			return &network.Data_Frame{
				Type:    network.Data_Message,
				Message: ackdata,
			}
		}
		return nil
	}
	return nil
}

func (s *Agent) ackhandler(ack []interface{}) interface{} {
	if ack == nil {
		return nil
	}
	switch len(ack) {
	case 1:
		return ack[0]
	case 2:
		s.sess.UserData = []interface{}{ack[1]}
		return ack[0]
	default:
	}
	return nil
}
