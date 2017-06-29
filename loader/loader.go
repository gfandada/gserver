// 封装了装载器，统一加载配置文件
// FIXME 暂时以固定的路径加载
package loader

import (
	"github.com/gfandada/gserver/util"

	"github.com/koding/multiconfig"
)

type Loader struct {
}

var timer *util.TimerCfg
var msg *util.MsgCfg

/***************************实现imodule接口*******************************/

func (loader *Loader) OnInit() {
	// 加载timer配置
	timer = new(util.TimerCfg)
	multTimer := multiconfig.NewWithPath("../../../cfg/timer.json")
	multTimer.MustLoad(timer)
	// 加载message配置
	msg = new(util.MsgCfg)
	multMsg := multiconfig.NewWithPath("../../../cfg/message.json")
	multMsg.MustLoad(msg)
}

func (loader *Loader) OnDestroy() {

}

func (loader *Loader) Run(ChClose chan bool) {
	//	// 每日定时
	//	for index, value := range timer.Daily {
	//		fmt.Println(index, value.Hour, value.Sec)
	//	}
	//	// 每周定时
	//	for index, value := range timer.Weekly {
	//		fmt.Println(index, value.Hour, value.Sec)
	//	}
	//	// 每月定时
	//	for index, value := range timer.Monthly {
	//		fmt.Println(index, value.Hour, value.Sec)
	//	}
	// 同步消息
	//	for _, value := range msg.Sync {
	//		rawMsgType(value.m)
	//	}
}

// 解析生成真实的消息类型
func rawMsgType(msgType string) {
	//	switch msgType {
	//	case "login_req":
	//		return protomsg.LoginReq
	//	case "login_ack":
	//		return protomsg.LoginAck
	//	}
}
