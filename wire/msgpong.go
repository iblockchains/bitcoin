package wire

// MsgPong 响应 ping 消息的 pong 消息, 主要用来确认连接是否仍然有效
// implements the Message interface and represents a bitcoin pong
// message which is used primarily to confirm that a connection is still valid
// in response to a bitcoin ping message (MsgPing).
//
// This message was not added until protocol versions AFTER BIP0031Version.
type MsgPong struct {
	Nonce uint64
}
