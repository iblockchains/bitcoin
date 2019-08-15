package wire

// MsgVersion 实现了 Message 接口,代表比特币版本消息.
// 一旦同外部连接时,节点用它来广播自己.远程节点用此信息同自身的进行协商.
// 然后,伴随着verack消息,远程节点会回复一个包含协商之后版本值的版本消息.
// implements the Message interface and represents a bitcoin version
// message.  It is used for a peer to advertise itself as soon as an outbound
// connection is made.  The remote peer then uses this information along with
// its own to negotiate.  The remote peer must then respond with a version
// message of its own containing the negotiated values followed by a verack
// message (MsgVerAck).  This exchange must take place before any further
// communication is allowed to proceed.
type MsgVersion struct {
}
