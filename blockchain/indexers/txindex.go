package indexers

import (
	"github.com/iblockchains/bitcoin/database"
)

// TxIndex implements a transaction by hash index.  That is to say, it supports
// querying all transactions by their hash.
type TxIndex struct {
	db         database.DB
	curBlockID uint32
}

// NewTxIndex 生成一个索引实例,用来生成区块链中的所有交易的哈希值
// 同每个区块,在区块中的位置和交易大小的映射
// 它实现了 Indexer 接口
//
// returns a new instance of an indexer that is used to create a
// mapping of the hashes of all transactions in the blockchain to the respective
// block, location within the block, and size of the transaction.
//
// It implements the Indexer interface which plugs into the IndexManager that in
// turn is used by the blockchain package.  This allows the index to be
// seamlessly maintained along with the chain.
func NewTxIndex(db database.DB) *TxIndex {
	return &TxIndex{db: db}
}
