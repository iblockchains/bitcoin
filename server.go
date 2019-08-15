package main

import (
	"fmt"
	"sync"

	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/database"
)

// server 提供一个比特币服务端同其它节点进行通信.
// provides a bitcoin server for handling communications to and from
// bitcoin peers.
type server struct {
	wg sync.WaitGroup
}

// newServer 返回一个监听特定地址的比特币服务节点.用于接收对等节点的连接.
// returns a new btcd server configured to listen on addr for the
// bitcoin network type specified by chainParams.  Use start to begin accepting
// connections from peers.
func newServer(listenAddrs, agentBlacklist, agentWhitelist []string,
	db database.DB, chainParams *chaincfg.Params,
	interrupt <-chan struct{}) (*server, error) {
	fmt.Println("待:newServer")
	return nil, nil
}

// Start 启动服务
// begins accepting connections from peers.
func (s *server) Start() {
	fmt.Println("待:server.Start")
}

// Stop 通过暂停所有同其它节点的连接和主接听器优雅的关闭服务器
// gracefully shuts down the server by stopping and disconnecting all
// peers and the main listener.
func (s *server) Stop() error {
	fmt.Println("待:server.Stop")
	return nil
}

// WaitForShutdown 阻塞直到主监听器和节点的处理程序都停止.
// blocks until the main listener and peer handlers are stopped.
func (s *server) WaitForShutdown() {
	fmt.Println("待:WaitForShutdown")
	// s.wg.Wait()
}
