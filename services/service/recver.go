package service

import (
	"fmt"
	"io"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
	Services "github.com/gfandada/gserver/services"
)

type recver struct {
	stream network.Service_StreamServer
	sess   *Session
	in     chan *network.Data_Frame
}

func (re *recver) run() {
	defer func() {
		close(re.in)
		Services.GetHandler(Services.CLOSE_CONNECT)([]interface{}{re.sess})
		logger.Debug("user %d close connection", re.sess.UserId)
	}()
	for {
		data, err := re.stream.Recv()
		if err == io.EOF { // 正常流关闭
			logger.Debug(fmt.Sprintf("Stream recver %v io.EOF", re))
			return
		}
		if err != nil { // 流错误
			logger.Error(fmt.Sprintf("Stream recver %v error %v", re, err))
			return
		}
		select {
		case re.in <- data:
		case <-re.sess.Die:
			return
		}
	}
}

func startRecver(stream network.Service_StreamServer, sess *Session) chan *network.Data_Frame {
	ch := make(chan *network.Data_Frame, 1)
	re := &recver{
		stream: stream,
		sess:   sess,
		in:     ch,
	}
	go re.run()
	return ch
}
