package indexers

import (
	"github.com/iblockchains/bitcoin/database"
)

// Manager 定义一了个可以管理多个可选索引器的索引管理器,
// 同时实现了 blockchain.IndexManager 接口,这样它就能
// 无缝的传入普通区块链处理过程.
// defines an index manager that manages multiple optional indexes and
// implements the blockchain.IndexManager interface so it can be seamlessly
// plugged into normal chain processing.
type Manager struct {
	db             database.DB
	enabledIndexes []Indexer
}

// NewManager 返回一个索引管理器
// returns a new index manager with the provided indexes enabled.
//
// The manager returned satisfies the blockchain.IndexManager interface and thus
// cleanly plugs into the normal blockchain processing path.
func NewManager(db database.DB, enabledIndexes []Indexer) *Manager {
	return &Manager{
		db:             db,
		enabledIndexes: enabledIndexes,
	}
}
