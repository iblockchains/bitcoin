package indexers

import (
	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/database"
)

// CfIndex 实现了通过区块哈希索引 committed filter
// implements a committed filter (cf) by hash index.
type CfIndex struct {
	db          database.DB
	chainParams *chaincfg.Params
}

// NewCfIndex 返回一索引实例,用于创建区块哈希值到对应 committed filters 的映射.
// returns a new instance of an indexer that is used to create a
// mapping of the hashes of all blocks in the blockchain to their respective
// committed filters.
func NewCfIndex(db database.DB, chainParams *chaincfg.Params) *CfIndex {
	return &CfIndex{db: db, chainParams: chainParams}
}
