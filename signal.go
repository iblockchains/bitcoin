package main

import "fmt"

// shutdownRequestChannel 用于接收从任一子系统的打断信号
// is used to initiate shutdown from one of the
// subsystems using the same code paths as when an interrupt signal is received.
var shutdownRequestChannel = make(chan struct{})

// interruptListener 接收从操作系统发出像(Ctrl+C)和来自shutdownRequestChannel的关闭请求.
// 当接收以上任一信号时就返回一个通道.
// listens for OS Signals such as SIGINT (Ctrl+C) and shutdown
// requests from shutdownRequestChannel.  It returns a channel that is closed
// when either signal is received.
func interruptListener() <-chan struct{} {
	fmt.Println("待:interruptListener")
	c := make(chan struct{})
	return c
}

// interruptRequested returns true when the channel returned by
// interruptListener was closed.  This simplifies early shutdown slightly since
// the caller can just use an if statement instead of a select.
func interruptRequested(interrupted <-chan struct{}) bool {
	// select {
	// case <-interrupted:
	// 	fmt.Println("interrupted ...........")
	// 	return true
	// default:
	// }
	fmt.Println("待:interruptRequested")
	return false
}
