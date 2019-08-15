package main

import "github.com/iblockchains/bitcoin/chaincfg"

// activeNetParams 一个特别指向当前比特币网络参数的指针
// activeNetParams is a pointer to the parameters specific to the
// currently active bitcoin network.
var activeNetParams = &mainNetParams

// 用于分组各种网络的参数，如主网络和测试网络
// is used to group parameters for various networks such as the main
// network and test networks
type params struct {
	*chaincfg.Params
	rpcPort string
}

var mainNetParams = params{
	Params:  &chaincfg.MainNetParams,
	rpcPort: "8334",
}
