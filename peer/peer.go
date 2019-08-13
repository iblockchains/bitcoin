package peer

import "net"

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
