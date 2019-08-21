package mining

import (
	"fmt"

	"github.com/iblockchains/bitcoin/blockchain"
	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/txscript"
)

// BlkTmplGenerator provides a type that can be used to generate block templates
// based on a given mining policy and source of transactions to choose from.
// It also houses additional state required in order to ensure the templates
// are built on top of the current best chain and adhere to the consensus rules.
type BlkTmplGenerator struct {
}

// TxSource 待打包进新区块的交易源
// 所有方法并发安全
// represents a source of transactions to consider for inclusion in
// new blocks.
//
// The interface contract requires that all of these methods are safe for
// concurrent access with respect to the source.
type TxSource interface {
}

// NewBlkTmplGenerator 返回一个使用指定策略的区块模块生成器
// returns a new block template generator for the given
// policy using transactions from the provided transaction source.
//
// The additional state-related fields are required in order to ensure the
// templates are built on top of the current best chain and adhere to the
// consensus rules.
func NewBlkTmplGenerator(policy *Policy, params *chaincfg.Params,
	txSource TxSource, chain *blockchain.BlockChain,
	timeSource blockchain.MedianTimeSource,
	sigCache *txscript.SigCache,
	hashCache *txscript.HashCache) *BlkTmplGenerator {
	fmt.Println("Unfinished:mining.BlkTmpGenerator")
	return nil
}
