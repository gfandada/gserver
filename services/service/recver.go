package service

import (
	"fmt"
	"io"

	"github.com/gfandada/gserver/logger"
	"github.com/gfandada/gserver/network"
)

type recver struct {
	stream network.Service_StreamServer
	die    <-chan struct{}
	in     chan *network.Data_Frame
}

func (re *recver) run() {
	defer func() {
		close(re.in)
	}()
	for {
		data, err := re.stream.Recv()
		if err == io.EOF { // 正常流关闭
			logger.Debug("Stream recver io.EOF")
			return
		}
		if err != nil { // 流错误
			logger.Error(fmt.Sprintf("Stream recver error %v", err))
			return
		}
		select {
		case re.in <- data:
		case <-re.die:
			return
		}
	}
}

func startRecver(stream network.Service_StreamServer, die chan struct{}) chan *network.Data_Frame {
	ch := make(chan *network.Data_Frame, 1)
	re := &recver{
		stream: stream,
		die:    die,
		in:     ch,
	}
	go re.run()
	logger.Debug("Stream recver run")
	return ch
}
