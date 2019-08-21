package indexers

import (
	"fmt"

	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/database"
)

// AddrIndex 根据地址查询交易
// implements a transaction by address index.  That is to say, it
// supports querying all transactions that reference a given address because
// they are either crediting or debiting the address.  The returned transactions
// are ordered according to their order of appearance in the blockchain.  In
// other words, first by block height and then by offset inside the block.
//
// In addition, support is provided for a memory-only index of unconfirmed
// transactions such as those which are kept in the memory pool before inclusion
// in a block.
type AddrIndex struct{}

// NewAddrIndex 返回一个索引实例,用于创建区块链中所有地址到对应交易的映射
// returns a new instance of an indexer that is used to create a
// mapping of all addresses in the blockchain to the respective transactions
// that involve them.
func NewAddrIndex(db database.DB, chainParams *chaincfg.Params) *AddrIndex {
	fmt.Println("Unfinished:indexers.NewAddrIndex")
	return nil
}
