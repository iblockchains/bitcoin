package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/iblockchains/bitcoin/txscript"

	"github.com/iblockchains/bitcoin/blockchain"
	"github.com/iblockchains/bitcoin/connmgr"
	"github.com/iblockchains/bitcoin/mempool"
	"github.com/iblockchains/bitcoin/mining"
	"github.com/iblockchains/bitcoin/mining/cpuminer"
	"github.com/iblockchains/bitcoin/netsync"

	"github.com/iblockchains/bitcoin/addrmgr"
	"github.com/iblockchains/bitcoin/blockchain/indexers"
	"github.com/iblockchains/bitcoin/wire"

	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/database"
)

const (
	// 描述了服务器所支持的默认服务
	// defaultServices describes the default services that are supported by
	// the server.
	// defaultServices = wire.SFNodeNetwork | wire.SFNodeBloom |
	// 	wire.SFNodeWitness | wire.SFNodeCF
	// defaultTargetOutbound 默认对外连接节点的个数
	// is the default number of outbound peers to target.
	defaultTargetOutbound = 8
)

// server 提供一个比特币服务端同其它节点进行通信.
// provides a bitcoin server for handling communications to and from
// bitcoin peers.
type server struct {
	chainParams *chaincfg.Params
	connManager *connmgr.ConnManager
	addrManager *addrmgr.AddrManager
	sigCache    *txscript.SigCache
	hashCache   *txscript.HashCache
	rpcServer   *rpcServer
	chain       *blockchain.BlockChain
	syncManager *netsync.SyncManager
	txMemPool   *mempool.TxPool
	cpuMiner    *cpuminer.CPUMiner
	services    wire.ServiceFlag
	timeSource  blockchain.MedianTimeSource
	nat         NAT
	db          database.DB
	wg          sync.WaitGroup

	// 可选索引.初始创建之后不再变化,所以无需并发安全
	// The following fields are used for optional indexes.  They will be nil
	// if the associated index is not enabled.  These fields are set during
	// initial creation of the server and never changed afterwards, so they
	// do not need to be protected for concurrent access.
	txIndex   *indexers.TxIndex
	addrIndex *indexers.AddrIndex
	cfIndex   *indexers.CfIndex

	// 跟踪 mempool 中的交易数量
	// The fee estimator keeps track of how long transactions are left in
	// the mempool before they are mined into blocks.
	feeEstimator *mempool.FeeEstimator
}

// simpleAddr 实现了net.Addr 接口
// implements the net.Addr interface with two struct fields
type simpleAddr struct {
	net, addr string
}

// String returns the address.
//
// This is part of the net.Addr interface.
func (a simpleAddr) String() string {
	return a.addr
}

// Network returns the network.
//
// This is part of the net.Addr interface.
func (a simpleAddr) Network() string {
	return a.net
}
func getServicesServerSupport() wire.ServiceFlag {
	var services wire.ServiceFlag
	// services = defaultServices
	// if cfg.NoPeerBloomFilters {
	// &^ 为 bit clear 运算符,services对应 wire.SFNodeBloom为1的位清0,
	// 如,0110 &^ 1011 = 0100,1011 &^ 1101 = 0010
	// 	services &^= wire.SFNodeBloom
	// }
	// if cfg.NoCFilters {
	// 	services &^= wire.SFNodeCF
	// }
	fmt.Println("Unfinished: getServicesServerSupport")
	return services
}

// initListeners 初始化配置的网络监听器和添加外部地址到地址管理器
// 返回监听器和NAT接口
// initializes the configured net listeners and adds any bound
// addresses to the address manager. Returns the listeners and a NAT interface,
// which is non-nil if UPnP is in use.
func initListeners(amgr *addrmgr.AddrManager, listenAddrs []string, services wire.ServiceFlag) ([]net.Listener, NAT, error) {
	// 在配置地址上监听TCP连接
	// Listen for TCP connections at the configured addresses
	netAddrs, err := parseListeners(listenAddrs)
	if err != nil {
		return nil, nil, err
	}
	listeners := make([]net.Listener, 0, len(netAddrs))
	for _, addr := range netAddrs {
		listener, err := net.Listen(addr.Network(), addr.String())
		if err != nil {
			srvrLog.Warnf("Can't listen on %s: %v", addr, err)
			continue
		}
		listeners = append(listeners, listener)
	}
	var nat NAT
	if len(cfg.ExternalIPs) != 0 {
		fmt.Println("Unfinished:server.initListeners")
	} else {
		if cfg.Upnp {
			fmt.Println("Unfinished:server.initListeners")
		}
		// 向地址管理器添加用于广播的绑定地址
		// Add bound addresses to address manager to be advertised to peers.
		for _, listener := range listeners {
			addr := listener.Addr().String()
			err := addLocalAddress(amgr, addr, services)
			if err != nil {
				amgrLog.Warnf("Skipping bound address %s: %v", addr, err)
			}
		}
	}
	return listeners, nat, nil
}

// parseListeners 判定每个地址是否是IPv4和IPv6并返回正确的地址.
//
// determines whether each listen address is IPv4 and IPv6 and
// returns a slice of appropriate net.Addrs to listen on with TCP. It also
// properly detects addresses which apply to "all interfaces" and adds the
// address as both IPv4 and IPv6.
func parseListeners(addrs []string) ([]net.Addr, error) {
	netAddrs := make([]net.Addr, 0, len(addrs)*2)
	for _, addr := range addrs {
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		// Empty host or host of * on plan9 is both IPv4 and IPv6.
		if host == "" || (host == "*" && runtime.GOOS == "plan9") {
			netAddrs = append(netAddrs, simpleAddr{net: "tcp4", addr: addr})
			netAddrs = append(netAddrs, simpleAddr{net: "tcp6", addr: addr})
			continue
		}
		// Strip IPv6 zone id if present since net.ParseIP does not
		// handle it.
		zoneIndex := strings.LastIndex(host, "%")
		if zoneIndex > 0 {
			host = host[:zoneIndex]
		}
		// Parse the IP
		ip := net.ParseIP(host)
		if ip == nil {
			return nil, fmt.Errorf("'%s' is not a valid IP address", host)
		}
		// To4 如果不是IPv4地址就返回nil
		if ip.To4() == nil {
			netAddrs = append(netAddrs, simpleAddr{net: "tcp6", addr: addr})
		} else {
			netAddrs = append(netAddrs, simpleAddr{net: "tcp4", addr: addr})
		}
	}
	return netAddrs, nil
}

// addLocalAddress 向地址管理器添加当前节点正在监听的地址,这样它就能传递给其它节点
// adds an address that this node is listening on to the
// address manager so that it may be relayed to peers.
func addLocalAddress(addrMgr *addrmgr.AddrManager, addr string, services wire.ServiceFlag) error {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}
	port, err := strconv.ParseUint(portStr, 10, 16) // 10 十进制,16类型 int16
	if err != nil {
		return err
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsUnspecified() { // 未指定IP
		// 如果绑定了未指定的地址,广播所有本地接口 (默认情况是未指定 0.0.0.0)
		// If bound to unspecified address, advertise all local interfaces
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			ifaceIP, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			// If bound to 0.0.0.0, do not add IPv6 interfaces and if bound to
			// ::, do not add IPv4 interfaces
			if (ip.To4() == nil) != (ifaceIP.To4() == nil) {
				continue
			}
			netAddr := wire.NewNetAddressIPPort(ifaceIP, uint16(port), services)
			addrMgr.AddLocalAddress(netAddr, addrmgr.BoundPrio)
		}
	} else {
		// netAddr, err := addrMgr.HostToNetAddress(host, uint(port), services)
		// if err != nil {
		// 	return err
		// }
		log.Panic("Unfinished:addrMgrHostToNetAddress")
		// addrMgr.AddLocalAddress(netAddr, addrmgr.BoundPrio)
	}
	return nil
}

// newServer 返回一个监听特定地址的比特币服务节点.用于接收对等节点的连接.
// returns a new btcd server configured to listen on addr for the
// bitcoin network type specified by chainParams.  Use start to begin accepting
// connections from peers.
func newServer(listenAddrs, agentBlacklist, agentWhitelist []string,
	db database.DB, chainParams *chaincfg.Params,
	interrupt <-chan struct{}) (*server, error) {
	services := getServicesServerSupport()
	// -----------------地址管理---------------------
	amgr := addrmgr.New(cfg.DataDir, btcdLookup)
	var listeners []net.Listener
	var nat NAT
	if !cfg.DisableListen {
		var err error
		listeners, nat, err = initListeners(amgr, listenAddrs, services)
		if err != nil {
			return nil, err
		}
		if len(listeners) == 0 {
			return nil, errors.New("no valid listen address")
		}
	}
	if len(agentBlacklist) > 0 {
		srvrLog.Infof("User-agent blacklist %s", agentBlacklist)
	}
	if len(agentWhitelist) > 0 {
		srvrLog.Infof("User-agent whitelist %s", agentWhitelist)
	}
	s := server{
		addrManager: amgr,
		nat:         nat,
	}
	// ------------- 索引管理 ------------------
	// 如何需要就新建交易和地址索引
	//
	// 注意: 要先新建交易索引,因为地址索引会用到交易索引的数据.
	//
	// Create the transaction and address indexes if needed.
	//
	// CAUTION: the txindex needs to be first in the indexes array because
	// the addrindex uses data from the txindex during catchup.  If the
	// addrindex is run first, it may not have the transactions from the
	// current block indexed.
	var indexes []indexers.Indexer
	if cfg.TxIndex || cfg.AddrIndex { // 添加交易索引
		// Enable transaction index if address index is enabled since it
		// requires it.
		if !cfg.TxIndex {
			indxLog.Infof("Transaction index enabled because it " +
				"is required by the address index")
			cfg.TxIndex = true
		} else {
			indxLog.Info("Transaction index is enabled")
		}
		s.txIndex = indexers.NewTxIndex(db)
		indexes = append(indexes, s.txIndex)
	}
	if cfg.AddrIndex { // 添加地址索引
		indxLog.Info("Address index is enabled")
		s.addrIndex = indexers.NewAddrIndex(db, chainParams)
		indexes = append(indexes, s.addrIndex)
	}
	if !cfg.NoCFilters {
		indxLog.Info("Committed filter index is enabled")
		s.cfIndex = indexers.NewCfIndex(db, chainParams)
		indexes = append(indexes, s.cfIndex)
	}
	// 任一可选的索引被启用,就新建索引管理器
	// Create an index manager if any of the optional indexes are enabled.
	var indexManager blockchain.IndexManager
	if len(indexes) > 0 {
		indexManager = indexers.NewManager(db, indexes)
	}
	// 将给定的检查点和默认的合并,除非他们被禁用
	// Merge given checkpoints with the default ones unless they are disabled
	var checkpoints []chaincfg.Checkpoint
	if !cfg.DisableCheckpoints {
		// checkpoints = mergeCheckpoints(s.chainParams.Checkpoints, cfg.addCheckpoints)
		fmt.Println("Unfinished:server.mergeCheckpoints")
	}
	// ---------- 新建区块链 -------------------
	// 使用合适的配置新建区块链实例
	// Create a new block chain instance with the appropriate configuration.
	var err error
	s.chain, err = blockchain.New(&blockchain.Config{
		DB:           s.db,
		Interrupt:    interrupt,
		ChainParams:  s.chainParams,
		Checkpoints:  checkpoints,
		TimeSource:   s.timeSource,
		SigCache:     s.sigCache,
		IndexManager: indexManager,
		HashCache:    s.hashCache,
	})
	if err != nil {
		return nil, err
	}
	// -------------- FeeEstimator -----------------
	// 在数据库中查询 FeeEstimator 的状态. 找不到或无法加载就新建
	// Search for a FeeEstimator state in the database. If none can be found
	// or if it cannot be loaded, create a new one.
	// db.Update(func(tx database.Tx) error {
	// 	fmt.Println("Unfinished: find  feeEstimationData")
	// 	return nil
	// })
	fmt.Println("Unfinished:server.FeeEstimator")
	// 启动新的FeeEstimator
	// If no feeEstimator has been found, or if the one that has been found
	// is behind somehow, create a new one and start over.
	if s.feeEstimator == nil || s.feeEstimator.LastKnownHeight() != s.chain.BestSnapshot().Height {
		fmt.Println("Unfinished: create new feeEstimator")
	}
	// -------------- 交易池 ------------------
	txC := mempool.Config{}
	s.txMemPool = mempool.New(&txC)
	// -------------- 同步管理 ---------------
	s.syncManager, err = netsync.New(&netsync.Config{})
	if err != nil {
		return nil, err
	}
	// --------------- mining -----------------
	// 创建采矿政策和区块模板生成器
	// 注意: CPU矿工依赖于交易池,因此交易池的创建必须先于CPU矿工
	//
	// Create the mining policy and block template generator based on the
	// configuration options.
	//
	// NOTE: The CPU miner relies on the mempool, so the mempool has to be
	// created before calling the function to create the CPU miner.
	policy := mining.Policy{}
	blockTemplateGenerator := mining.NewBlkTmplGenerator(&policy,
		s.chainParams, s.txMemPool, s.chain, s.timeSource,
		s.sigCache, s.hashCache)
	s.cpuMiner = cpuminer.New(&cpuminer.Config{
		BlockTemplateGenerator: blockTemplateGenerator,
	})

	// 当不在只能连接模式下运行时，仅设置一个函数以返回要连接的新地址。
	// 模拟网络一直处于只能连接的模式下, 因为它仅用于连接指定的节点,
	// 和积极避免广播和连接到发现的节点,以防止成为一个公共测试网络.
	//
	// Only setup a function to return new addresses to connect to when
	// not running in connect-only mode.  The simulation network is always
	// in connect-only mode since it is only intended to connect to
	// specified peers and actively avoid advertising and connecting to
	// discovered peers in order to prevent it from becoming a public test
	// network.
	var newAddressFunc func() (net.Addr, error)
	if !cfg.SimNet && len(cfg.ConnectPeers) == 0 { // 不是模拟网络并且未连接节点
		newAddressFunc = func() (net.Addr, error) {
			// for tries := 0; tries < 100; tries++ {
			// 	addr := s.addrManager.GetAddress()
			// 	if addr == nil {
			// 		break
			// 	}
			// 	// 地址不会是无效的,或是本地的地址,或是不可路由的,
			// 	// 因为 addrmanager 添加时会拒绝这些类型的地址.
			// 	// 务必检查一下同一组没有相同的地址,这样我们就不会
			// 	// 以牺牲他人为代价连接到相同的网络片断
			// 	// Address will not be invalid, local or unroutable
			// 	// because addrmanager rejects those on addition.
			// 	// Just check that we don't already have an address
			// 	// in the same group so that we are not connecting
			// 	// to the same network segment at the expense of
			// 	// others.
			// 	key:=addrmgr.GroupKey(addr.NetAddress())
			// }
			fmt.Println("Unfinished: newAddressFunc")
			return nil, errors.New("no valid connect address")
		}
	}
	// ---------------- 连接管理 Connection Manager ------------
	// 新建一个连接管理器
	// Create a connection manager
	targetOutbound := defaultTargetOutbound
	if cfg.MaxPeers < targetOutbound {
		targetOutbound = cfg.MaxPeers
	}
	cmgr, err := connmgr.New(&connmgr.Config{
		GetNewAddress: newAddressFunc,
	})
	if err != nil {
		return nil, err
	}
	s.connManager = cmgr
	// 启动持久化的节点
	// Start up persistent peers
	permanentPeers := cfg.ConnectPeers
	if len(permanentPeers) == 0 {
		permanentPeers = cfg.AddPeers
	}
	for _, addr := range permanentPeers {
		netAddr, err := addrStringToNetAddr(addr)
		if err != nil {
			return nil, err
		}
		go s.connManager.Connect(&connmgr.ConnReq{
			Addr:      netAddr,
			Permanent: true,
		})
	}
	// --------------- RPC -----------
	if !cfg.DisableRPC {
		// 为配置好的RPC设置地址和TLS配置监听器
		// Setup listeners for the configured RPC listen addresses and
		// TLS settings.
		// rpcListeners, err := setupRPCListeners()
		// if err != nil {
		// 	return nil, err
		// }
		// if len(rpcListeners) == 0 {
		// 	return nil, errors.New("RPCS: No valid listen address")
		// }
		// s.rpcServer, err = newRPCServer(&rpcserverConfig{})
		// if err != nil {
		// 	return nil, err
		// }
		// go func() {
		// 	<-s.rpcServer.RequestedProcessShutdown()
		// 	shutdownRequestChannel <- struct{}{}
		// }()
		fmt.Println("Unfinished:server.RPC")
	}
	fmt.Println("Unfinished:server.newServer  End")
	return &s, nil
}

// Start 启动服务
// begins accepting connections from peers.
func (s *server) Start() {
	fmt.Println("Unfinished:server.Start")
}

// Stop 通过暂停所有同其它节点的连接和主接听器优雅的关闭服务器
// gracefully shuts down the server by stopping and disconnecting all
// peers and the main listener.
func (s *server) Stop() error {
	fmt.Println("Unfinished:server.Stop")
	return nil
}

// WaitForShutdown 阻塞直到主监听器和节点的处理程序都停止.
// blocks until the main listener and peer handlers are stopped.
func (s *server) WaitForShutdown() {
	fmt.Println("Unfinished:WaitForShutdown")
	// s.wg.Wait()
}

// mergeCheckpoints returns two slices of checkpoints merged into one slice
// such that the checkpoints are sorted by height.  In the case the additional
// checkpoints contain a checkpoint with the same height as a checkpoint in the
// default checkpoints, the additional checkpoint will take precedence and
// overwrite the default one.
func mergeCheckpoints(defaultCheckpoints, additional []chaincfg.Checkpoint) []chaincfg.Checkpoint {
	fmt.Println("Unfinished:server.mergeCheckpoints")
	return nil
}
func addrStringToNetAddr(addr string) (net.Addr, error) {
	fmt.Println("Unfinished:server.addrStringToNetAddr")
	return nil, nil
}
func setupRPCListeners() ([]net.Listener, error) {
	fmt.Println("Unfinished:server.setupRPCLiseners")
	return nil, nil
}
