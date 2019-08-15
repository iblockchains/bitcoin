package chaincfg

// Params 根据参数定义一个比特币网络.
// defines a Bitcoin network by its parameters.  These parameters may be
// used by Bitcoin applications to differentiate networks as well as addresses
// and keys for one network from those intended for use on another network.
type Params struct {
}

// MainNetParams 为比特币主网络定义网络参数
// defines the network parameters for the main Bitcoin network.
var MainNetParams = Params{}
