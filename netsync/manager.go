package netsync

import "fmt"

// SyncManager 用于节点间与块相关的消息通信.
// 使用Start()启动一个协程.启动后会选择节点下载和同步区块.
// is used to communicate block related messages with peers. The
// SyncManager is started as by executing Start() in a goroutine. Once started,
// it selects peers to sync from and starts the initial block download. Once the
// chain is in sync, the SyncManager handles incoming block and header
// notifications and relays announcements of new blocks to peers.
type SyncManager struct{}

func New(config *Config) (*SyncManager, error) {
	fmt.Println("Unfinished:netsync.New")
	return nil, nil
}
