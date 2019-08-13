package wire

import "io"

// MessageEncoding 表示 wire 消息的编码格式
// represents the wire message encoding format to be used
type MessageEncoding uint32

const (
	// BaseEncoding 按照比特币连接协议指定的格式转码所有的消息
	// encodes all messages in the default format specified
	// for the Bitcoin wire protocol.
	BaseEncoding MessageEncoding = 1 << iota // 1
	// WitnessEncoding 按照特定比特币连接协议转码所有除交易外的所有信息.
	// 对于交易信息会使用BIP0144定义的转码格式
	// encodes all messages other than transaction messages
	// using the default Bitcoin wire protocol specification. For transaction
	// messages, the new encoding format detailed in BIP0144 will be used.
	WitnessEncoding // 2
)

// Message 比特币消息方法接口.
// is an interface that describes a bitcoin message.  A type that
// implements Message has complete control over the representation of its data
// and may therefore contain additional or fewer fields than those which
// are used directly in the protocol encoded message.
type Message interface {
	BtcDecode(io.Reader, uint32, MessageEncoding) error
	BtcEncode(io.Reader, uint32, MessageEncoding) error
	Command() string
	MaxPayloadLength(uint32) uint32
}
