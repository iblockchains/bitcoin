package addrmgr

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"github.com/iblockchains/bitcoin/wire"
)

// AddrManager 为比特币网络上可缓存的节点提供了一个并发安全的地址管理器
// provides a concurrency safe address manager for caching potential
// peers on the bitcoin network.
type AddrManager struct {
	lamtx          sync.Mutex
	localAddresses map[string]*localAddress
}
type localAddress struct {
	na    *wire.NetAddress
	score AddressPriority
}

// AddressPriority 地址类型
// type is used to describe the hierarchy of local address
// discovery methods.
type AddressPriority int

const (
	// InterfacePrio 表示地址属于本地接口
	// signifies the address is on a local interface
	InterfacePrio AddressPriority = iota
	// BoundPrio 表示地址有明确的限制
	// signifies the address has been explicitly bounded to.
	BoundPrio
	// UpnpPrio 表示地址源于UPnP
	// signifies the address was obtained from UPnP.
	UpnpPrio
	// HTTPPrio 表示地址源于外部HTTP服务
	// signifies the address was obtained from an external HTTP service.
	HTTPPrio
	// ManualPrio 表示地址由--externalip提供
	// signifies the address was provided by --externalip.
	ManualPrio
)

// New 返回一个比特币地址管理器
// 使用 Start 开始处理异步地址更新
// returns a new bitcoin address manager.
// Use Start to begin processing asynchronous address updates.
func New(dataDir string, lookupFunc func(string) ([]net.IP, error)) *AddrManager {
	am := AddrManager{
		localAddresses: make(map[string]*localAddress),
	}
	fmt.Println("Unfinished:addrmgr.New")
	am.reset()
	return &am
}

// reset 重置地址管理器
// resets the address manager by reinitialising the random source
// and allocating fresh empty bucket storage.
func (a *AddrManager) reset() {
	fmt.Println("Unfinished:addrmgr.reset")
}

// AddLocalAddress 根据优先级将地址添加到本地广播已知地址列表
// adds na to the list of known local addresses to advertise
// with the given priority.
func (a *AddrManager) AddLocalAddress(na *wire.NetAddress, priority AddressPriority) error {
	if !IsRoutable(na) { // 不可路由
		return fmt.Errorf("address %s is not routable", na.IP)
	}
	a.lamtx.Lock()
	defer a.lamtx.Unlock()

	key := NetAddressKey(na) //生成网络地址的key
	la, ok := a.localAddresses[key]
	if !ok || la.score < priority {
		if ok {
			la.score = priority + 1
		} else {
			a.localAddresses[key] = &localAddress{
				na:    na,
				score: priority,
			}
		}
	}
	return nil
}

// ipString ip转换成字符串格式.
// returns a string for the ip from the provided NetAddress. If the
// ip is in the range used for Tor addresses then it will be transformed into
// the relevant .onion address.
func ipString(na *wire.NetAddress) string {
	if IsOnionCatTor(na) {
		log.Panic("Unfinished: addrmanager.ipString")
	}
	return na.IP.String()
}

// NetAddressKey returns a string key in the form of ip:port for IPv4 addresses
// or [ip]:port for IPv6 addresses.
func NetAddressKey(na *wire.NetAddress) string {
	port := strconv.FormatUint(uint64(na.Port), 10)
	return net.JoinHostPort(ipString(na), port)
}
