// 网络对象
package entity

import (
	Network "github.com/gfandada/gserver/network"
	Service "github.com/gfandada/gserver/services/service"
)

type GameClient struct {
	clientid int32
}

func (game *GameClient) Post(msg Network.RawMessage) {
	if game.clientid > 0 {
		Service.Send(game.clientid, msg)
	}
}
