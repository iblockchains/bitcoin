package chaincfg

import "github.com/btcsuite/btcd/wire"

// Checkpoint 标识区块链中已知的好点.
// 使用 checkpoints 在区块初始下载中可以对旧区块做一些优化
// 也能防止摘取旧的区块
type Checkpoint struct {
}

// Params 根据参数定义一个比特币网络.
// defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
	// Name 网络名
	// defines a human-readable identifier for the network.
	Name string
	// Net defines a
	Net wire.BitcoinNet
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string
	// Checkpoints 从旧到新排序
	// ordered from oldest to newest.
	Checkpoints []Checkpoint
}

// MainNetParams 为比特币主网络定义网络参数
// defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{
	Name:        "mainnet",
	Net:         wire.MainNet,
	DefaultPort: "8333",
}
