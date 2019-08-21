package wire

// ServiceFlag 确认一个比特币节点支持的服务
// identifies services supported by a bitcoin peer.
type ServiceFlag uint64

// BitcoinNet 比特币网络类型
// represents which bitcoin network a message belongs to.
type BitcoinNet uint32

const (
	// MainNet 代表比特币主网络
	// represents the main bitcoin network.
	MainNet BitcoinNet = 0xd9b4bef9
	// TestNet 递归测试网络
	// represents the regression test network.
	TestNet BitcoinNet = 0xdab5bffa
	// TestNet3 版本3测试网络
	// represents the test network (version 3).
	TestNet3 BitcoinNet = 0x0709110b

	// SimNet 模拟测试网络.
	// represents the simulation test network
	SimNet BitcoinNet = 0x12141c16
)
