package gateway

import (
	"time"

	"github.com/gfandada/gserver/network"
)

const (
	SESS_AUTHFAILED  = 0x4 // 认证失败
	SESS_AUTHSUCCEED = 0x8 // 认证成功
)

type Session struct {
	MQ                chan network.Data_Frame      // 返回给客户端的异步消息
	UserId            int32                        // 玩家ID
	Stream            network.Service_StreamClient // 后端游戏服数据流
	Die               chan struct{}                // 会话关闭信号
	Flag              int32                        // 会话标记
	ConnectTime       time.Time                    // 链接建立时间
	PacketTime        time.Time                    // 当前包的到达时间
	LastPacketTime    time.Time                    // 上一个包到达时间
	PacketCount       uint32                       // 对收到的包进行计数
	PacketCountOneMin int                          // 每分钟的包统计，用于RPM判断
}
