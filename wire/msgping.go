package wire

// MsgPing 比特币 ping 消息接口
// 对于BIP0031Version及之前的版本,它主要是用来确认连接是否仍然有效.
// 传输错误通常会被理解为连接关闭并移除对应的节点.
// 对于BIP0031Version之后的版本,在返回的PONG信息中,它包含一个标识符,以确定网络时间
// implements the Message interface and represents a bitcoin ping
// message.
//
// For versions BIP0031Version and earlier, it is used primarily to confirm
// that a connection is still valid.  A transmission error is typically
// interpreted as a closed connection and that the peer should be removed.
// For versions AFTER BIP0031Version it contains an identifier which can be
// returned in the pong message to determine network timing.
//
// The payload for this message just consists of a nonce used for identifying
// it later.
type MsgPing struct {
	// Unique value associated with message that is used to identify
	// specific ping message.
	Nonce uint64
}
