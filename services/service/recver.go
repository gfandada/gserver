package service

import (
	"io"

	"../../network"
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
			return
		}
		if err != nil { // 流错误
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
	return ch
}
