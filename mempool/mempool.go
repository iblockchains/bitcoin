package mempool

import "fmt"

// Config is a descriptor containing the memory pool configuration.
type Config struct {
}

// TxPool 用于缓存未打包进区块的交易.
// 并发安全
// is used as a source of transactions that need to be mined into blocks
// and relayed to other peers.  It is safe for concurrent access from multiple
// peers.
type TxPool struct{}

// New 返回一个用于验证和储存独立交易的缓存池,直到他们被打包进区块
// returns a new memory pool for validating and storing standalone
// transactions until they are mined into a block.
func New(cfg *Config) *TxPool {
	fmt.Println("Unfinished:mempool.New")
	return nil
}
