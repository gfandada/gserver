package util

// 消息配置
type MsgCfg struct {
	Sync  []MsgInfo
	Async []MsgInfo
}

// 双向消息
type MsgInfo struct {
	ReqMsgId   int
	ReqMsgData string
	AckMsgId   int
	AckMsgData string
}
