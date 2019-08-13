package peer

import (
	"fmt"
	"net"

	"github.com/iblockchains/bitcoin/wire"
)

// MessageListeners 定义了节点的消息监听器的回调函数指针。
// 在节点初始化时，任何未设置具体回调的监听器会被忽略。
// 多个消息监听器是串联执行的，上一个执行结束前会阻塞下一个的执行
// 注意：除非另有说明，否则这些监听器不能直接调用节点实例上的任何阻塞回调函数，因为在回调完成之前，输入处理程序goroutine会阻塞。这样做将导致死锁。
// defines callback function pointers to invoke with message listeners for a peer
// Any listener which is not set to a concrete callback during peer initialization is ignored.
// Execution of multiple message listeners occurs serially, so one callback blocks the execution of the next.
// NOTE: Unless otherwise documented, these listeners must NOT directly call any
// blocking calls (such as WaitForShutdown) on the peer instance since the input
// handler goroutine blocks until the callback has completed.  Doing so will
// result in a deadlock.
type MessageListeners struct {
}

// Config 节点配置信息
// is the struct to hold configuration options useful to Peer.
type Config struct {
}

// newNetAddress 分解传入地址的IP和端口，并新建一个比特币地址
// newNetAddress attempts to extract the IP address and port from the passed
// net.Addr interface and create a bitcoin NetAddress structure using that
// information.
func newNetAddress(addr net.Addr, services wire.ServiceFlag) (*wire.NetAddress, error) {
}
func minUint32(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

// outMsg 它用于存储要发送的消息，以及在消息发送时发出信号的通道
// is used to house a message to be sent along with a channel to signal
// when the message has been sent
type outMsg struct {
}

// StatsSnap 某个时间节点数据分析快照
// is a snapshot of peer stats at a point in time
type StatsSnap struct {
}

// type HashFunc func() (hash *chainhash.Hash, height int32, err error)

// HashFunc 用于获取最新区块信息的回调函数
// is a function which returns a block hash, height and error
// It is used as a callback to get newest block details.
type HashFunc func()

// AddrFunc 获取一个地址并返回一个相关地址
// is a func which takes an address and returns a related address.
type AddrFunc func(remoteAddr *wire.NetAddress) *wire.NetAddress

// HostToNetAddrFunc 根据host、port和服务返回netaddress
// is a func which takes an address and returns a related address
type HostToNetAddrFunc func(host string, port uint16, services wire.ServiceFlag) (*wire.NetAddress, error)

// 注意：节点的数据流被分为三个goroutines来处理。
// 传入信息：通过 inHandler 线程来读取，并且通常会分发到自己的处理程序。
// 对传入信息的相关数据，如区块、交易和库存，会被相关的程序处理。
// 传出数据：通过两个线程 queueHandler  和 outHandler 来处理数据流。
// queueHandler 是一种对外部实体进行排队的方式

// NOTE: The overall data flow of a peer is split into 3 goroutines.  Inbound
// messages are read via the inHandler goroutine and generally dispatched to
// their own handler.  For inbound data-related messages such as blocks,
// transactions, and inventory, the data is handled by the corresponding
// message handlers.  The data flow for outbound messages is split into 2
// goroutines, queueHandler and outHandler.  The first, queueHandler, is used
// as a way for external entities to queue messages, by way of the QueueMessage
// function, quickly regardless of whether the peer is currently sending or not.
// It acts as the traffic cop between the external world and the actual
// goroutine which writes to the network socket.

// Peer provides a basic concurrent safe bitcoin peer for handling bitcoin
// communications via the peer-to-peer protocol.  It provides full duplex
// reading and writing, automatic handling of the initial handshake process,
// querying of usage statistics and other information about the remote peer such
// as its address, user agent, and protocol version, output message queuing,
// inventory trickling, and the ability to dynamically register and unregister
// callbacks for handling bitcoin protocol messages.
//
// Outbound messages are typically queued via QueueMessage or QueueInventory.
// QueueMessage is intended for all messages, including responses to data such
// as blocks and transactions.  QueueInventory, on the other hand, is only
// intended for relaying inventory as it employs a trickling mechanism to batch
// the inventory together.  However, some helper functions for pushing messages
// of specific types that typically require common special handling are
// provided as a convenience.

// Peer 为基于点对点协议的比特币通信提供一个基本的并发安全的比特币节点。
// 它提供双向的读写操作，自动处理初始握手过程，查询远程节点的使用统计和
// 其它信息，像它的 address、user agent 和 protocol 版本等

// Peer provides a basic concurrent safe bitcoin peer for handling bitcoin
// communications via the peer-to-peer protocol.  It provides full duplex
// reading and writing, automatic handling of the initial handshake process,
// querying of usage statistics and other information about the remote peer such
// as its address, user agent, and protocol version, output message queuing,
// inventory trickling, and the ability to dynamically register and unregister
// callbacks for handling bitcoin protocol messages.
type Peer struct {
	addr    string
	inbound bool
}

// String 返回人类可读的节点的地址和方向性(传出或传入)
// 并发安全
// String returns the peer's address and directionality as a human-readable
// string.
//
// This function is safe for concurrent access.
func (p *Peer) String() string {
	return fmt.Sprintf("%s (%s)", p.addr, directionString(p.inbound))
}

// UpdateLastBlockHeight 更新节点最后一个已知区块的高度
// 并发安全
// UpdateLastBlockHeight updates the last known block for the peer.
//
// This function is safe for concurrent access.
func (p *Peer) UpdateLastBlockHeight(newHeight int32) {
	fmt.Println("UpdateLastBlockHeight")
}

// UpdateLastAnnouncedBlock 更新节点最后广播区块的哈希值
// 并发安全
// UpdateLastAnnouncedBlock updates meta-data about the last block hash this
// peer is known to have announced.
//
// This function is safe for concurrent access.
func (p *Peer) UpdateLastAnnouncedBlock() {
	fmt.Println("UpdateLastAnnouncedBlock")
}

// AddKnownInventory 添加已知传递的清单（inventory）到已知清单缓存
// 并发安全
// AddKnownInventory adds the passed inventory to the cache of known inventory
// for the peer.
//
// This function is safe for concurrent access.
func (p *Peer) AddKnownInventory() {
	fmt.Println(" AddKnownInventory")
}

// StatsSnapshot 返回当前节点的统计数据
// 并发安全
// StatsSnapshot returns a snapshot of the current peer flags and statistics.
//
// This function is safe for concurrent access.
func (p *Peer) StatsSnapshot() *StatsSnap {
	fmt.Println("StatsSnapshot")
	return nil
}

// ID returns the peer id.
//
// This function is safe for concurrent access.
func (p *Peer) ID() int32 {
	fmt.Println("返回节点ID")
	return 0
}

// NA 返回的是节点的网络地址
// NA returns the peer network address.
//
// This function is safe for concurrent access.
func (p *Peer) NA() *wire.NetAddress {
	fmt.Println("NA")
	return nil
}

// Addr returns the peer address.
//
// This function is safe for concurrent access.
func (p *Peer) Addr() string {
	// The address doesn't change after initialization, therefore it is not
	// protected by a mutex.
	return p.addr
}

// Inbound returns whether the peer is inbound.
//
// This function is safe for concurrent access.
func (p *Peer) Inbound() bool {
	return p.inbound
}

// Services 返回远程节点的服务标志
// Services returns the services flag of the remote peer.
//
// This function is safe for concurrent access.
func (p *Peer) Services() wire.ServiceFlag {
	fmt.Println("Services")
	return 0
}

// UserAgent returns the user agent of the remote peer.
//
// This function is safe for concurrent access.
func (p *Peer) UserAgent() string {
	fmt.Println("UserAgent")
	return ""
}

// LastAnnouncedBlock returns the last announced block of the remote peer.
//
// This function is safe for concurrent access.
func (p *Peer) LastAnnouncedBlock() {
	fmt.Println("LastAnnouncedBlock")
}

// PushAddrMsg 发送地址消息到邻近节点.
// 此函数在通过 QueueMessage 手动发送消息时非常有用.
// 当地址太多时会自动限制地址数量为上限并随机选择地址.
// 它返回实际发送的地址,如果提供的地址切片中没有条目，则不会发送消息.
// sends an addr message to the connected peer using the provided
// addresses.  This function is useful over manually sending the message via
// QueueMessage since it automatically limits the addresses to the maximum
// number allowed by the message and randomizes the chosen addresses when there
// are too many.  It returns the addresses that were actually sent and no
// message will be sent if there are no entries in the provided addresses slice.
//
// This function is safe for concurrent access.
func (p *Peer) PushAddrMsg(addresses []*wire.NetAddress) ([]*wire.NetAddress, error) {
	fmt.Println("PushAddrMsg")
	return nil, nil
}

// PushGetBlocksMsg 发送getblocks消息
// sends a getblocks message for the provided block locator
// and stop hash.  It will ignore back-to-back duplicate requests.
//
// This function is safe for concurrent access.
func (p *Peer) PushGetBlocksMsg() error {
	fmt.Println("PushGetBlocksMsg")
	return nil
}

// PushGetHeadersMsg 发送 getheaders 消息
// sends a getblocks message for the provided block locator
// and stop hash.  It will ignore back-to-back duplicate requests.
//
// This function is safe for concurrent access.
func (p *Peer) PushGetHeadersMsg() error {
	fmt.Println("PushGetHeadersMsg")
	return nil
}

// PushRejectMsg 发送 reject 消息
// sends a reject message for the provided command, reject code,
// reject reason, and hash.  The hash will only be used when the command is a tx
// or block and should be nil in other cases.  The wait parameter will cause the
// function to block until the reject message has actually been sent.
//
// This function is safe for concurrent access.
func (p *Peer) PushRejectMsg() {
	fmt.Println("PushRejectMsg")
}

// handlePingMsg 当节点接收到一个ping比特币消息时被调用.
// 对于协议版本大于BIP0031Version的客户端会回复一个pong消息.
// is invoked when a peer receives a ping bitcoin message.  For
// recent clients (protocol version > BIP0031Version), it replies with a pong
// message.  For older clients, it does nothing and anything other than failure
// is considered a successful ping.
func (p *Peer) handlePingMsg(msg *wire.MsgPing) {
	fmt.Println("handlePingMsg")
}

// handlePongMsg 当节点收到一个 pong 消息时被调用,它更新 ping 时间统计.
// is invoked when a peer receives a pong bitcoin message.  It
// updates the ping statistics as required for recent clients (protocol
// version > BIP0031Version).  There is no effect for older clients or when a
// ping was not previously sent.
func (p *Peer) handlePongMsg(msg *wire.MsgPong) {
	fmt.Println("handlePongMsg")
}

// readMessage 从节点读取比特币消息并纪录日志
// reads the next bitcoin message from the peer with logging.
func (p *Peer) readMessage(encoding wire.MessageEncoding) (wire.Message, []byte, error) {
	fmt.Println("readMessage")
	return nil, nil, nil
}

// writeMessage 发送比特币消息给节点并纪录日志
// sends a bitcoin message to the peer with logging.
func (p *Peer) writeMessage(msg wire.Message, enc wire.MessageEncoding) error {
	fmt.Println("writeMessage")
	return nil
}

// stallHandler 处理节点的速检测.考虑到回调所花费的时间,
// 这需要跟踪预期的响应并为其分配截止时间
// handles stall detection for the peer.  This entails keeping
// track of expected responses and assigning them deadlines while accounting for
// the time spent in callbacks.  It must be run as a goroutine.
func (p *Peer) stallHandler() {
	fmt.Println("stallHandler")
}

// inHandler 处理所有传入节点的消息
// handles all incoming messages for the peer.
// It must be run as a goroutine.
func (p *Peer) inHandler() {
	fmt.Println("inHandler")
}

// queueHandler 节点传出数据队列化.对各种输入源来说它就像一个
// 路由器,这样我们就能保证, 我们发送消息时不会被服务端和节点阻塞.
// 最终所有数据都会传递给outHandler
// handles the queuing of outgoing data for the peer. This runs as
// a muxer for various sources of input so we can ensure that server and peer
// handlers will not block on us sending a message.  That data is then passed on
// to outHandler to be actually written.
func (p *Peer) queueHandler() {
	fmt.Println("queueHandler")
}

// outHandler 处理所有节点传出数据.它必须做为一个goroutine运行.
// 它使用缓存通道序列化输出消息,同时允许发送方异步运行
// handles all outgoing messages for the peer.  It must be run as a
// goroutine.  It uses a buffered channel to serialize output messages while
// allowing the sender to continue running asynchronously.
func (p *Peer) outHandler() {
	fmt.Println("outHandler")
}

// pingHandler 定期 pings 节点.
// periodically pings the peer.  It must be run as a goroutine.
func (p *Peer) pingHandler() {
	fmt.Println("pingHandler")
}

// QueueMessage adds the passed bitcoin message to the peer send queue.
//
// This function is safe for concurrent access.
func (p *Peer) QueueMessage(msg wire.Message, doneChan chan<- struct{}) {
	p.QueueMessageWithEncoding(msg, doneChan, wire.BaseEncoding)
}

// QueueMessageWithEncoding 添加传递的消息到节点的发送队列.
// adds the passed bitcoin message to the peer send
// queue. This function is identical to QueueMessage, however it allows the
// caller to specify the wire encoding type that should be used when
// encoding/decoding blocks and transactions.
//
// This function is safe for concurrent access.
func (p *Peer) QueueMessageWithEncoding(msg wire.Message, doneChan chan<- struct{},
	encoding wire.MessageEncoding) {
	fmt.Println("QueueMessageWithEncoding")
}
