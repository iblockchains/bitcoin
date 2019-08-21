<h1>搭建基本的启动服务框架 start a simple server</h1>

* [第一步](#1)

  * [主函数 btcd](#101)

  * [日志 log](#102)

  * [配置 config](#103)

  * [信号 signal](#104)

  * [升级 upgrade](#105)

  * [服务器 server](#106)

  * [版本信息 version](#107)

  * [参数](#108)

* [第二步 服务器架构和地址管理](#2)
  * [修改服务 update server](#201)
  * [地址管理 addrmgr](#202)
  * [索引管理 indexers](#203)
  * [区块链 blockchain](#204)
  * [数据库 database](#205)
  * [交易池 mempool](#206)
  * [同步管理 netsync](#207)
  * [挖矿 mining](#208)
  * [连接管理 Connection Manager](#209)
  * [RPC](#210)
<h3 id="1">第一步</h3>

<h5>客户端启动流程</h5>
<img src="https://github.com/iblockchains/bitcoin/blob/master/img/bitcoin-001-btcd-startup.png">

<h4 id="101">主函数 main</h4>
<h5>btcd.go:</h5>

```
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/iblockchains/bitcoin/database"
	"github.com/iblockchains/bitcoin/limits"
)

const (
	// blockDbNamePrefix 是区块数据库名的前缀.
	// is the prefix for the block database name.  The
	// database type is appended to this value to form the full block
	// database name.
	blockDbNamePrefix = "blocks"
)

var (
	cfg *config
)

// winServiceMain 只有在Windows上会被调用.它检测出btcd何时会以服务形式运行,
// 并做出反应.
// is only invoked on Windows.  It detects when btcd is running
// as a service and reacts accordingly.
var winServiceMain func() (bool, error)

// btcdMain 是btcd真正的主程序. 有必要解决这样的事实,当os.Exit()被调用之前,延迟函数不要执行.
// serverChan 是一个可选的参数,它被服务代码用于在服务器启动后向其发出通知，以便在服务控制管
// 理器请求时能够优雅地停止服务.
// is the real main function for btcd.  It is necessary to work around
// the fact that deferred functions do not run when os.Exit() is called.  The
// optional serverChan parameter is mainly used by the service code to be
// notified with the server once it is setup so it can gracefully stop it when
// requested from the service control manager.
func btcdMain(serverChan chan<- *server) error {
	// 加载配置并分析命令行.此函数还启动日志并对其进行相应的配置
	// Load configuration and parse command line.  This function also
	// initializes logging and configures it accordingly.
	tcfg, _, err := loadConfig()
	if err != nil {
		return err
	}
	cfg = tcfg
	defer func() { // 最后关闭日志
		if logRotator != nil {
			logRotator.Close()
		}
	}()
	// 获取这样一个通道,当无论从操作系统(Ctrl+C)或从
	// 子系统像RPC服务端发出关闭信号后,这个通道就将关闭
	// Get a channel that will be closed when a shutdown signal has been
	// triggered either from an OS signal such as SIGINT (Ctrl+C) or from
	// another subsystem such as the RPC server.
	interrupt := interruptListener()
	defer btcdLog.Info("完全关闭.Shutdown complete")
	// 启动时显示版本
	// Show version at startup
	// btcdLog.Info("Version %s", version())

	// 根据需要启动http分析服务器
	// Enable http profiling server if requested.
	if cfg.Profile != "" {
		fmt.Println("Unfinished:Enable http profiling server")
	}
	// Write cpu profile if requested.
	if cfg.CPUProfile != "" {
		fmt.Println("Unfinished:Write cpu profile if requested")
	}
	// 版本升级
	// Perform upgrades to btcd as new versions require it
	if err := doUpgrades(); err != nil {
		btcdLog.Errorf("%v", err)
		return err
	}
	// 如果触发了打断信号就返回
	if interruptRequested(interrupt) {
		return nil
	}
	// 加载区块数据库
	// Load the block database
	db, err := loadBlockDB()
	if err != nil {
		btcdLog.Errorf("%v", err)
		return err
	}
	defer func() {
		btcdLog.Infof("Unfinished: Gracefully shutting down the database...")
		if db != nil {
			db.Close()
		}
	}()
	// 如果需要就丢弃索引,并退出
	//
	// 注意: 这里的顺序很重要,由于依赖关系,丢弃 tx 索引时也丢弃了地址索引
	// Drop indexes and exit if requested.
	//
	// NOTE: The order is important here because dropping the tx index also
	// drops the address index since it relies on it.
	if cfg.DropAddrIndex { // 丢弃地址索引
		fmt.Println("Unfinished:DropAddrIndex")
		return nil
	}
	if cfg.DropTxIndex {
		fmt.Println("Unfinished:DropTxIndex")
		return nil
	}
	if cfg.DropCfIndex {
		fmt.Println("Unfinished:DropCfIndex")
		return nil
	}
	// 新建服务器并启动
	// Create server and start it.
	server, err := newServer(cfg.Listeners, cfg.AgentBlacklist,
		cfg.AgentWhitelist, db, activeNetParams.Params, interrupt)
	if err != nil {
		btcdLog.Errorf("Unable to start server on %v:%v", cfg.Listeners, err)
		return err
	}

	defer func() {
		btcdLog.Infof("Gracefully shutting down the server...")
		server.Stop()
		server.WaitForShutdown()
		srvrLog.Infof("Server shutdown complete")
	}()
	server.Start()
	if serverChan != nil {
		serverChan <- server
	}
	// 阻塞直到收到关闭信号
	// Wait until the interrupt signal is received from an OS signal or
	// shutdown is requested through one of the subsystems such as the RPC
	// server.
	// <-interrupt
	return nil
}
func main() {
	//使用所有的处理器
	// Use all processor cores.
	runtime.GOMAXPROCS(runtime.NumCPU())
	// 区块和交易的处理会产生突发性的资源分配.
	// 这限制了垃圾收集器在突发期间的过度分配.
	// Block and transaction processing can cause bursty allocations.  This
	// limits the garbage collector from excessively overallocating during
	// bursts.  This value was arrived at with the help of profiling live
	// usage.
	debug.SetGCPercent(10)
	// 设置某些限制
	// Up some limits.
	if err := limits.SetLimits(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set limits:%v\n", err)
		os.Exit(1)
	}
	// Call serviceMain on Windows to handle running as a service.  When
	// the return isService flag is true, exit now since we ran as a
	// service.  Otherwise, just fall through to normal operation.
	if runtime.GOOS == "windows" {
		isService, err := winServiceMain()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if isService {
			os.Exit(0)
		}
	}
	// Work around defer not working after os.Exit()
	if err := btcdMain(nil); err != nil {
		os.Exit(1)
	}
}
func loadBlockDB() (database.DB, error) {
	fmt.Println("Unfinished:loadBlockDB")
	return nil, nil
}

```

<h4 id="102">日志 log</h4>
<h5>日志在整个客户端的各个角落随处可见,先将它的基本框架和必要功能先写上.</h5>
<h5>log.go:</h5>
<h5>首先,需要使用go get拉取两个包github.com/btcsuite/btclog和github.com/jrick/logrotate/rotator</h5>

```
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/btcsuite/btclog"
	"github.com/jrick/logrotate/rotator"
)

// logWriter 实现了io.Writer接口可以同时将日志打印到控制台和输出到log rotator
// implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized log rotator.
type logWriter struct{}

// Write logWriter引用io.Writer接口的具体实现
func (logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	logRotator.Write(p)
	return len(p), nil
}

// 每个子系统的纪录器.只创建一个后端纪录器,所有的子系统将基于此创建各自的纪录器
// Loggers per subsystem.  A single backend logger is created and all subsytem
// loggers created from it will write to the backend.  When adding new
// subsystems, add the subsystem logger variable here and to the
// subsystemLoggers map.
//
// Loggers can not be used before the log rotator has been initialized with a
// log file.  This must be performed early during application startup by calling
// initLogRotator.
var (
	// backendLog 日志纪录后端用于创建子系统的日志纪录器.
	// is the logging backend used to create all subsystem loggers.
	backendLog = btclog.NewBackend(logWriter{})
	// logRotator 是日志输出中的一个.它能从文件读取日志并将日志
	// 写入文件,当文件太大时它会压缩和截短文件.
	// is one of the logging outputs.
	logRotator *rotator.Rotator
	btcdLog    = backendLog.Logger("BTCD") // 客户端日志
	srvrLog    = backendLog.Logger("SRVR") // 服务器日志
  indxLog    = backendLog.Logger("INDX") // 索引日志
  amgrLog    = backendLog.Logger("AMGR") // 地址管理器
)

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func initLogRotator(logFile string) {
	// fmt.Printf("完:initLogRotator:%s\n", logFile)
	logDir, _ := filepath.Split(logFile) //获得路径名(不包含文件名和其后缀在内)
	// fmt.Println(logDir)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory:%v\n", err)
		os.Exit(1)
	}
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}

	logRotator = r
}
```

<h4 id="103">配置 config</h4>
<h5>go get github.com/btcsuite/btcutil</h5>
<h5>新建config.go文件:</h5>

```
package main

import (
	"fmt"
	"path/filepath"

	"github.com/btcsuite/btcutil"
)

const (
	defaultLogFilename = "btcd.log"
	defaultLogDirname  = "logs"
)

var (
	defaultHomeDir = btcutil.AppDataDir("btcd", false)
	defaultLogDir  = filepath.Join(defaultHomeDir, defaultLogDirname)
)

// config btcd配置定义
// config defines the configuration options for btcd.
//
// See loadConfig for details on the configuration load process.
type config struct {
	ShowVersion          bool          `short:"V" long:"version" description:"Display version information and exit"`
	ConfigFile           string        `short:"C" long:"configfile" description:"Path to configuration file"`
	DataDir              string        `short:"b" long:"datadir" description:"Directory to store data"`
	LogDir               string        `long:"logdir" description:"Directory to log output."`
	AddPeers             []string      `short:"a" long:"addpeer" description:"Add a peer to connect with at startup"`
	ConnectPeers         []string      `long:"connect" description:"Connect only to the specified peers at startup"`
	DisableListen        bool          `long:"nolisten" description:"Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen"`
	Listeners            []string      `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 8333, testnet: 18333)"`
	MaxPeers             int           `long:"maxpeers" description:"Max number of inbound and outbound peers"`
	DisableBanning       bool          `long:"nobanning" description:"Disable banning of misbehaving peers"`
	BanDuration          time.Duration `long:"banduration" description:"How long to ban misbehaving peers.  Valid time units are {s, m, h}.  Minimum 1 second"`
	BanThreshold         uint32        `long:"banthreshold" description:"Maximum allowed ban score before disconnecting and banning misbehaving peers."`
	Whitelists           []string      `long:"whitelist" description:"Add an IP network or IP that will not be banned. (eg. 192.168.1.0/24 or ::1)"`
	AgentBlacklist       []string      `long:"agentblacklist" description:"A comma separated list of user-agent substrings which will cause btcd to reject any peers whose user-agent contains any of the blacklisted substrings."`
	AgentWhitelist       []string      `long:"agentwhitelist" description:"A comma separated list of user-agent substrings which will cause btcd to require all peers' user-agents to contain one of the whitelisted substrings. The blacklist is applied before the blacklist, and an empty whitelist will allow all agents that do not fail the blacklist."`
	RPCUser              string        `short:"u" long:"rpcuser" description:"Username for RPC connections"`
	RPCPass              string        `short:"P" long:"rpcpass" default-mask:"-" description:"Password for RPC connections"`
	RPCLimitUser         string        `long:"rpclimituser" description:"Username for limited RPC connections"`
	RPCLimitPass         string        `long:"rpclimitpass" default-mask:"-" description:"Password for limited RPC connections"`
	RPCListeners         []string      `long:"rpclisten" description:"Add an interface/port to listen for RPC connections (default port: 8334, testnet: 18334)"`
	RPCCert              string        `long:"rpccert" description:"File containing the certificate file"`
	RPCKey               string        `long:"rpckey" description:"File containing the certificate key"`
	RPCMaxClients        int           `long:"rpcmaxclients" description:"Max number of RPC clients for standard connections"`
	RPCMaxWebsockets     int           `long:"rpcmaxwebsockets" description:"Max number of RPC websocket connections"`
	RPCMaxConcurrentReqs int           `long:"rpcmaxconcurrentreqs" description:"Max number of concurrent RPC requests that may be processed concurrently"`
	RPCQuirks            bool          `long:"rpcquirks" description:"Mirror some JSON-RPC quirks of Bitcoin Core -- NOTE: Discouraged unless interoperability issues need to be worked around"`
	DisableRPC           bool          `long:"norpc" description:"Disable built-in RPC server -- NOTE: The RPC server is disabled by default if no rpcuser/rpcpass or rpclimituser/rpclimitpass is specified"`
	DisableTLS           bool          `long:"notls" description:"Disable TLS for the RPC server -- NOTE: This is only allowed if the RPC server is bound to localhost"`
	DisableDNSSeed       bool          `long:"nodnsseed" description:"Disable DNS seeding for peers"`
	ExternalIPs          []string      `long:"externalip" description:"Add an ip to the list of local addresses we claim to listen on to peers"`
	Proxy                string        `long:"proxy" description:"Connect via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	ProxyUser            string        `long:"proxyuser" description:"Username for proxy server"`
	ProxyPass            string        `long:"proxypass" default-mask:"-" description:"Password for proxy server"`
	OnionProxy           string        `long:"onion" description:"Connect to tor hidden services via SOCKS5 proxy (eg. 127.0.0.1:9050)"`
	OnionProxyUser       string        `long:"onionuser" description:"Username for onion proxy server"`
	OnionProxyPass       string        `long:"onionpass" default-mask:"-" description:"Password for onion proxy server"`
	NoOnion              bool          `long:"noonion" description:"Disable connecting to tor hidden services"`
	TorIsolation         bool          `long:"torisolation" description:"Enable Tor stream isolation by randomizing user credentials for each connection."`
	TestNet3             bool          `long:"testnet" description:"Use the test network"`
	RegressionTest       bool          `long:"regtest" description:"Use the regression test network"`
	SimNet               bool          `long:"simnet" description:"Use the simulation test network"`
	AddCheckpoints       []string      `long:"addcheckpoint" description:"Add a custom checkpoint.  Format: '<height>:<hash>'"`
	DisableCheckpoints   bool          `long:"nocheckpoints" description:"Disable built-in checkpoints.  Don't do this unless you know what you're doing."`
	DbType               string        `long:"dbtype" description:"Database backend to use for the Block Chain"`
	Profile              string        `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	CPUProfile           string        `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DebugLevel           string        `short:"d" long:"debuglevel" description:"Logging level for all subsystems {trace, debug, info, warn, error, critical} -- You may also specify <subsystem>=<level>,<subsystem2>=<level>,... to set the log level for individual subsystems -- Use show to list available subsystems"`
	Upnp                 bool          `long:"upnp" description:"Use UPnP to map our listening port outside of NAT"`
	MinRelayTxFee        float64       `long:"minrelaytxfee" description:"The minimum transaction fee in BTC/kB to be considered a non-zero fee."`
	FreeTxRelayLimit     float64       `long:"limitfreerelay" description:"Limit relay of transactions with no transaction fee to the given amount in thousands of bytes per minute"`
	NoRelayPriority      bool          `long:"norelaypriority" description:"Do not require free or low-fee transactions to have high priority for relaying"`
	TrickleInterval      time.Duration `long:"trickleinterval" description:"Minimum time between attempts to send new inventory to a connected peer"`
	MaxOrphanTxs         int           `long:"maxorphantx" description:"Max number of orphan transactions to keep in memory"`
	Generate             bool          `long:"generate" description:"Generate (mine) bitcoins using the CPU"`
	MiningAddrs          []string      `long:"miningaddr" description:"Add the specified payment address to the list of addresses to use for generated blocks -- At least one address is required if the generate option is set"`
	BlockMinSize         uint32        `long:"blockminsize" description:"Mininum block size in bytes to be used when creating a block"`
	BlockMaxSize         uint32        `long:"blockmaxsize" description:"Maximum block size in bytes to be used when creating a block"`
	BlockMinWeight       uint32        `long:"blockminweight" description:"Mininum block weight to be used when creating a block"`
	BlockMaxWeight       uint32        `long:"blockmaxweight" description:"Maximum block weight to be used when creating a block"`
	BlockPrioritySize    uint32        `long:"blockprioritysize" description:"Size in bytes for high-priority/low-fee transactions when creating a block"`
	UserAgentComments    []string      `long:"uacomment" description:"Comment to add to the user agent -- See BIP 14 for more information."`
	NoPeerBloomFilters   bool          `long:"nopeerbloomfilters" description:"Disable bloom filtering support"`
	NoCFilters           bool          `long:"nocfilters" description:"Disable committed filtering (CF) support"`
	DropCfIndex          bool          `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
	SigCacheMaxSize      uint          `long:"sigcachemaxsize" description:"The maximum number of entries in the signature verification cache"`
	BlocksOnly           bool          `long:"blocksonly" description:"Do not accept transactions from remote peers."`
	TxIndex              bool          `long:"txindex" description:"Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	DropTxIndex          bool          `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
	AddrIndex            bool          `long:"addrindex" description:"Maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
	DropAddrIndex        bool          `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
	RelayNonStd          bool          `long:"relaynonstd" description:"Relay non-standard transactions regardless of the default settings for the active network."`
	RejectNonStd         bool          `long:"rejectnonstd" description:"Reject non-standard transactions regardless of the default settings for the active network."`
	RejectReplacement    bool          `long:"rejectreplacement" description:"Reject transactions that attempt to replace existing transactions within the mempool through the Replace-By-Fee (RBF) signaling policy."`
	lookup               func(string) ([]net.IP, error)
	oniondial            func(string, string, time.Duration) (net.Conn, error)
	dial                 func(string, string, time.Duration) (net.Conn, error)
	addCheckpoints       []chaincfg.Checkpoint
	miningAddrs          []btcutil.Address
	minRelayTxFee        btcutil.Amount
	whitelists           []*net.IPNet
}

// loadConfig 从文件和命令行初始和解析配置.
// 配置过程如下:
// 1) 从健全的默认配置开始
// 2) 预解析命令行,检查是否存在可替代的配置文件
// 3) 加载配置文件并覆盖默认配置
// 4) 解析命令行CLI可选配置并覆盖之前配置
// initializes and parses the config using a config file and command
// line options.
//
// The configuration proceeds as follows:
// 	1) Start with a default config with sane settings
// 	2) Pre-parse the command line to check for an alternative config file
// 	3) Load configuration file overwriting defaults with any specified options
// 	4) Parse CLI options and overwrite/add any specified options
//
// The above results in btcd functioning properly without any config settings
// while still allowing the user to override settings with config files and
// command line options.  Command line options always take precedence.
func loadConfig() (*config, []string, error) {
	cfg := config{
		LogDir: defaultLogDir,
	}
	fmt.Println("Unfinished:loadConfig")
	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename))
	return &cfg, nil, nil
}
```

<h4 id="104">信号 signal</h4>
<h5>signal.go</h5>

```
package main

import "fmt"

// shutdownRequestChannel 用于接收从任一子系统的打断信号
// is used to initiate shutdown from one of the
// subsystems using the same code paths as when an interrupt signal is received.
var shutdownRequestChannel = make(chan struct{})

// interruptListener 接收从操作系统发出像(Ctrl+C)和来自shutdownRequestChannel的关闭请求.
// 当接收以上任一信号时就返回一个通道.
// listens for OS Signals such as SIGINT (Ctrl+C) and shutdown
// requests from shutdownRequestChannel.  It returns a channel that is closed
// when either signal is received.
func interruptListener() <-chan struct{} {
	fmt.Println("Unfinished:interruptListener")
	c := make(chan struct{})
	return c
}

// interruptRequested returns true when the channel returned by
// interruptListener was closed.  This simplifies early shutdown slightly since
// the caller can just use an if statement instead of a select.
func interruptRequested(interrupted <-chan struct{}) bool {
	// select {
	// case <-interrupted:
	// 	fmt.Println("interrupted ...........")
	// 	return true
	// default:
	// }
	fmt.Println("Unfinished:interruptRequested")
	return false
}
```

<h4 id="105">升级 upgrade</h4>
<h5>新建 upgrade.go</h5>

```
package main

import "fmt"

// doUpgrades performs upgrades to btcd as new versions require it.
func doUpgrades() error {
	fmt.Println("Unfinished:doUpgrades")
	return nil
}
```

<h4 id="106">服务器 server</h4>
<h5> server.go</h5>

```
package main

import (
	"fmt"
	"sync"

	"github.com/iblockchains/bitcoin/chaincfg"
	"github.com/iblockchains/bitcoin/database"
)

// server 提供一个比特币服务端同其它节点进行通信.
// provides a bitcoin server for handling communications to and from
// bitcoin peers.
type server struct {
	wg sync.WaitGroup
}

// newServer 返回一个监听特定地址的比特币服务节点.用于接收对等节点的连接.
// returns a new btcd server configured to listen on addr for the
// bitcoin network type specified by chainParams.  Use start to begin accepting
// connections from peers.
func newServer(listenAddrs, agentBlacklist, agentWhitelist []string,
	db database.DB, chainParams *chaincfg.Params,
	interrupt <-chan struct{}) (*server, error) {
	fmt.Println("Unfinished:newServer")
	return nil, nil
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

```

<h4 id="107">版本信息</h4>
<h5>新建 version.go</h5>

```
package main

import "fmt"

func version() string {
	fmt.Println("Unfinished:version")
	return "version"
}
```

<h4 id="108">app参数</h4>
<h5>新建params.go</h5>

```
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
```

<h4>区块链参数</h4>
<h5>新建 chaincfg/params.go</h5>

```
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
```

<h4>启动测试1 test</h4>

<h5>启动测试 go build and run app.</h5>

```
Unfinished:serviceMain
Unfinished:loadConfig
Unfinished:interruptListener
Unfinished:doUpgrades
Unfinished:interruptRequested
Unfinished:loadBlockDB
Unfinished:newServer
Unfinished:server.Start
2019-08-16 09:15:44.908 [INF] BTCD: Gracefully shutting down the server...
Unfinished:server.Stop
Unfinished:WaitForShutdown
2019-08-16 09:15:44.919 [INF] SRVR: Server shutdown complete
2019-08-16 09:15:44.920 [INF] BTCD: Unfinished: Gracefully shutting down the database...
2019-08-16 09:15:44.920 [INF] BTCD: 完全关闭.Shutdown complete
```

<h5>至此,基本框架搭建完成,下面要开始填充</h5>

<h3 id="2">第二步</h3>

<img src="https://github.com/iblockchains/bitcoin/blob/master/img/001-bitcoin-server-startup.png">

<h4 id="201">修改服务 update server</h4>
<h5>server.go:</h5>

```
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
	// -------------- 数据库 -----------------
	// 在数据库中查询 FeeEstimator 的状态. 找不到或无法加载就新建
	// Search for a FeeEstimator state in the database. If none can be found
	// or if it cannot be loaded, create a new one.
	// db.Update(func(tx database.Tx) error {
	// 	fmt.Println("Unfinished: find  feeEstimationData")
	// 	return nil
	// })
	fmt.Println("Unfinished:server.DB")
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

```


<h4 id="202">地址管理 addrmgr</h4>
<h5>config.go</h5>

```
type config struct {
	Listeners      []string `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 8333, testnet: 18333)"`
	AgentBlacklist []string `long:"agentblacklist" description:"A comma separated list of user-agent substrings which will cause btcd to reject any peers whose user-agent contains any of the blacklisted substrings."`
	AgentWhitelist []string `long:"agentwhitelist" description:"A comma separated list of user-agent substrings which will cause btcd to require all peers' user-agents to contain one of the whitelisted substrings. The blacklist is applied before the blacklist, and an empty whitelist will allow all agents that do not fail the blacklist."`
	LogDir         string   `long:"logdir" description:"Directory to log output."`
	Profile        string   `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	CPUProfile     string   `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DropAddrIndex  bool     `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
	DropTxIndex    bool     `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
	DropCfIndex    bool     `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
	lookup         func(string) ([]net.IP, error)
	DisableListen  bool `long:"nolisten" description:"Disable listening for incoming connections -- NOTE: Listening is automatically disabled if the --connect or --proxy options are used without also specifying listen interfaces via --listen"`
	TxIndex        bool `long:"txindex" description:"Maintain a full hash-based transaction index which makes all transactions available via the getrawtransaction RPC"`
	AddrIndex      bool `long:"addrindex" description:"Maintain a full address-based transaction index which makes the searchrawtransactions RPC available"`
}

func loadConfig() (*config, []string, error) {
	cfg := config{
		LogDir: defaultLogDir,
	}
	// --proxy or --connect without --listen disables listening.
	if (cfg.Proxy != "" || len(cfg.ConnectPeers) > 0) &&
		len(cfg.Listeners) == 0 {
		cfg.DisableListen = true
	}
	if len(cfg.Listeners) == 0 {
		cfg.Listeners = []string{
			net.JoinHostPort("", activeNetParams.DefaultPort),
		}
	}
	fmt.Println("Unfinished:loadConfig")
	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename))
	return &cfg, nil, nil
}
func btcdLookup(host string) ([]net.IP, error) {
	if strings.HasSuffix(host, ".onion") { // host 为洋葱头网络地址
		return nil, fmt.Errorf("attempt to resolve tor address %s", host)
	}
	return cfg.lookup(host)
}
```
<h5>upnp.go:</h5>

```
package main

// NAT 网络地址转换,是一个代表NAT遍历选项的接口,如UPNP或 NAT-PMP
// 它提供了查询和操作此遍历的方法，以允许访问服务。
//
// is an interface representing a NAT traversal options for example UPNP or
// NAT-PMP. It provides methods to query and manipulate this traversal to allow
// access to services.
type NAT interface {
}
```

<h5>addrmgr/addrmanager.go:</h5>

```
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

```
<h5>addrmgr/network.go</h5>

```
var (
	// rfc1918Nets 根据RFC1918定义了IPv4专用地址
	// specifies the IPv4 private address blocks as defined by
	// by RFC1918 (10.0.0.0/8, 172.16.0.0/12, and 192.168.0.0/16).
	rfc1918Nets = []net.IPNet{
		ipNet("10.0.0.0", 8, 32),
		ipNet("172.16.0.0", 12, 32),
		ipNet("192.168.0.0", 16, 32),
	}
	// rfc2544Net specifies the the IPv4 block as defined by RFC2544
	// (198.18.0.0/15)
	rfc2544Net = ipNet("198.18.0.0", 15, 32)

	// rfc3849Net specifies the IPv6 documentation address block as defined
	// by RFC3849 (2001:DB8::/32).
	rfc3849Net = ipNet("2001:DB8::", 32, 128)

	// rfc3927Net specifies the IPv4 auto configuration address block as
	// defined by RFC3927 (169.254.0.0/16).
	rfc3927Net = ipNet("169.254.0.0", 16, 32)

	// rfc3964Net specifies the IPv6 to IPv4 encapsulation address block as
	// defined by RFC3964 (2002::/16).
	rfc3964Net = ipNet("2002::", 16, 128)

	// rfc4193Net specifies the IPv6 unique local address block as defined
	// by RFC4193 (FC00::/7).
	rfc4193Net = ipNet("FC00::", 7, 128)

	// rfc4380Net specifies the IPv6 teredo tunneling over UDP address block
	// as defined by RFC4380 (2001::/32).
	rfc4380Net = ipNet("2001::", 32, 128)

	// rfc4843Net specifies the IPv6 ORCHID address block as defined by
	// RFC4843 (2001:10::/28).
	rfc4843Net = ipNet("2001:10::", 28, 128)

	// rfc4862Net specifies the IPv6 stateless address autoconfiguration
	// address block as defined by RFC4862 (FE80::/64).
	rfc4862Net = ipNet("FE80::", 64, 128)

	// rfc5737Net specifies the IPv4 documentation address blocks as defined
	// by RFC5737 (192.0.2.0/24, 198.51.100.0/24, 203.0.113.0/24)
	rfc5737Net = []net.IPNet{
		ipNet("192.0.2.0", 24, 32),
		ipNet("198.51.100.0", 24, 32),
		ipNet("203.0.113.0", 24, 32),
	}

	// rfc6052Net specifies the IPv6 well-known prefix address block as
	// defined by RFC6052 (64:FF9B::/96).
	rfc6052Net = ipNet("64:FF9B::", 96, 128)

	// rfc6145Net specifies the IPv6 to IPv4 translated address range as
	// defined by RFC6145 (::FFFF:0:0:0/96).
	rfc6145Net = ipNet("::FFFF:0:0:0", 96, 128)

	// rfc6598Net specifies the IPv4 block as defined by RFC6598 (100.64.0.0/10)
	rfc6598Net = ipNet("100.64.0.0", 10, 32)
	// onionCatNet 定义了一个用于支持洋葱头网络的IPv6地址块
	// defines the IPv6 address block used to support Tor.
	// bitcoind encodes a .onion address as a 16 byte number by decoding the
	// address prior to the .onion (i.e. the key hash) base32 into a ten
	// byte number. It then stores the first 6 bytes of the address as
	// 0xfd, 0x87, 0xd8, 0x7e, 0xeb, 0x43.
	//
	// This is the same range used by OnionCat, which is part part of the
	// RFC4193 unique local IPv6 range.
	//
	// In summary the format is:
	// { magic 6 bytes, 10 bytes base32 decode of key hash }
	onionCatNet = ipNet("fd87:d87e:eb43::", 48, 128)

	// zero4Net defines the IPv4 address block for address staring with 0
	// (0.0.0.0/8).
	zero4Net = ipNet("0.0.0.0", 8, 32)

	// heNet defines the Hurricane Electric IPv6 address block.
	heNet = ipNet("2001:470::", 32, 128)
)

// ipNet returns a net.IPNet struct given the passed IP address string, number
// of one bits to include at the start of the mask, and the total number of bits
// for the mask.
func ipNet(ip string, ones, bits int) net.IPNet {
	return net.IPNet{IP: net.ParseIP(ip), Mask: net.CIDRMask(ones, bits)}
}

// IsRFC1918 returns whether or not the passed address is part of the IPv4
// private network address space as defined by RFC1918 (10.0.0.0/8,
// 172.16.0.0/12, or 192.168.0.0/16)
func IsRFC1918(na *wire.NetAddress) bool {
	for _, rfc := range rfc1918Nets {
		if rfc.Contains(na.IP) {
			return true
		}
	}
	return false
}

// IsRFC2544 returns whether or not the passed address is part of the IPv4
// address space as defined by RFC2544 (198.18.0.0/15)
func IsRFC2544(na *wire.NetAddress) bool {
	return rfc2544Net.Contains(na.IP)
}

// IsRFC3849 returns whether or not the passed address is part of the IPv6
// documentation range as defined by RFC3849 (2001:DB8::/32).
func IsRFC3849(na *wire.NetAddress) bool {
	return rfc3849Net.Contains(na.IP)
}

// IsRFC3927 returns whether or not the passed address is part of the IPv4
// autoconfiguration range as defined by RFC3927 (169.254.0.0/16).
func IsRFC3927(na *wire.NetAddress) bool {
	return rfc3927Net.Contains(na.IP)
}

// IsRFC3964 returns whether or not the passed address is part of the IPv6 to
// IPv4 encapsulation range as defined by RFC3964 (2002::/16).
func IsRFC3964(na *wire.NetAddress) bool {
	return rfc3964Net.Contains(na.IP)
}

// IsRFC4193 returns whether or not the passed address is part of the IPv6
// unique local range as defined by RFC4193 (FC00::/7).
func IsRFC4193(na *wire.NetAddress) bool {
	return rfc4193Net.Contains(na.IP)
}

// IsRFC4380 returns whether or not the passed address is part of the IPv6
// teredo tunneling over UDP range as defined by RFC4380 (2001::/32).
func IsRFC4380(na *wire.NetAddress) bool {
	return rfc4380Net.Contains(na.IP)
}

// IsRFC4843 returns whether or not the passed address is part of the IPv6
// ORCHID range as defined by RFC4843 (2001:10::/28).
func IsRFC4843(na *wire.NetAddress) bool {
	return rfc4843Net.Contains(na.IP)
}

// IsRFC4862 returns whether or not the passed address is part of the IPv6
// stateless address autoconfiguration range as defined by RFC4862 (FE80::/64).
func IsRFC4862(na *wire.NetAddress) bool {
	return rfc4862Net.Contains(na.IP)
}

// IsRFC5737 returns whether or not the passed address is part of the IPv4
// documentation address space as defined by RFC5737 (192.0.2.0/24,
// 198.51.100.0/24, 203.0.113.0/24)
func IsRFC5737(na *wire.NetAddress) bool {
	for _, rfc := range rfc5737Net {
		if rfc.Contains(na.IP) {
			return true
		}
	}

	return false
}

// IsRFC6052 returns whether or not the passed address is part of the IPv6
// well-known prefix range as defined by RFC6052 (64:FF9B::/96).
func IsRFC6052(na *wire.NetAddress) bool {
	return rfc6052Net.Contains(na.IP)
}

// IsRFC6145 returns whether or not the passed address is part of the IPv6 to
// IPv4 translated address range as defined by RFC6145 (::FFFF:0:0:0/96).
func IsRFC6145(na *wire.NetAddress) bool {
	return rfc6145Net.Contains(na.IP)
}

// IsRFC6598 returns whether or not the passed address is part of the IPv4
// shared address space specified by RFC6598 (100.64.0.0/10)
func IsRFC6598(na *wire.NetAddress) bool {
	return rfc6598Net.Contains(na.IP)
}

// IsLocal 判断是否是本地地址
// returns whether or not the given address is a local address.
func IsLocal(na *wire.NetAddress) bool {
	// IsLoopback() 判断是否是一个回环地址,在windows 127.0.0.1 是一个回环地址
	return na.IP.IsLoopback() || zero4Net.Contains(na.IP)
}

// IsOnionCatTor returns whether or not the passed address is in the IPv6 range
// used by bitcoin to support Tor (fd87:d87e:eb43::/48).  Note that this range
// is the same range used by OnionCat, which is part of the RFC4193 unique local
// IPv6 range.
func IsOnionCatTor(na *wire.NetAddress) bool {
	return onionCatNet.Contains(na.IP)
}

// IsRoutable 判定当前地址在公共网络中是否可以路由.
// returns whether or not the passed address is routable over
// the public internet.  This is true as long as the address is valid and is not
// in any reserved ranges.
func IsRoutable(na *wire.NetAddress) bool {
	return IsValid(na) && !(IsRFC1918(na) || IsRFC2544(na) ||
		IsRFC3927(na) || IsRFC4862(na) || IsRFC3849(na) ||
		IsRFC4843(na) || IsRFC5737(na) || IsRFC6598(na) ||
		IsLocal(na) || (IsRFC4193(na) && !IsOnionCatTor(na)))
}

// IsValid 判断地址是否有效
// returns whether or not the passed address is valid.  The address is
// considered invalid under the following circumstances:
// IPv4: It is either a zero or all bits set address.
// IPv6: It is either a zero or RFC3849 documentation address.
func IsValid(na *wire.NetAddress) bool {
	return na.IP != nil && !(na.IP.IsUnspecified()) ||
		na.IP.Equal(net.IPv4bcast)
}

```
<h5>addrmgr/network_test.go</h5>

```
func TestIsRoutable(t *testing.T) {
	fmt.Println("==============IP IsRoutable==========")
	localAddress := "127.0.0.1"
	localNetAddress := wire.NetAddress{IP: net.ParseIP(localAddress)}
	fmt.Printf("%s isValid: %v\n", localAddress, IsValid(&localNetAddress))
	fmt.Printf("%s isLocal: %v\n", localAddress, IsLocal(&localNetAddress))
	fmt.Printf("%s IsRoutable: %v\n", localAddress, IsRoutable(&localNetAddress))
	fmt.Println("---------------------------------")
	fmt.Println()
}
```

<h5>$ cd addrmgr</h5>
<h5>$ go test -v network_test.go network.go</h5>
<h5>wire/netaddress.go</h5>

```
// NetAddress 定义网络中一个对等节点的信息，包括支持的服务、IP和端口
// NetAddress defines information about a peer on the network including the time
// it was last seen, the services it supports, its IP address, and port.
type NetAddress struct {
	// 地址纪录时间
	// Last time the address was seen.
	Timestamp time.Time
	// 地址所支持的服务
	// Bitfield which identifies the services supported by the address
	Services ServiceFlag
	// 节点IP
	// IP address of the peer
	IP net.IP
	// 节点端口
	// Port the peer is using.
	Port uint16
}

// NewNetAddressIPPort returns a new NetAddress using the provided IP, port, and
// supported services with defaults for the remaining fields.
func NewNetAddressIPPort(ip net.IP, port uint16, services ServiceFlag) *NetAddress {
	return NewNetAddressTimestamp(time.Now(), services, ip, port)
}

// NewNetAddressTimestamp returns a new NetAddress using the provided
// timestamp, IP, port, and supported services. The timestamp is rounded to
// single second precision.
func NewNetAddressTimestamp(
	timestamp time.Time, services ServiceFlag, ip net.IP, port uint16) *NetAddress {
	na := NetAddress{
		Timestamp: time.Unix(timestamp.Unix(), 0),
		Services:  services,
		IP:        ip,
		Port:      port,
	}
	return &na
}

```

<h4 id="203">索引管理 indexers</h4>

<h5>blockchain/indexers/common.go:</h5>

```
// Indexer 为一个受索引管理器管理的索引器提供通用的接口
// provides a generic interface for an indexer that is managed by an
// index manager such as the Manager type provided by this package.
type Indexer interface{}
```

<h5>blockchain/indexers/txindex.go:</h5>

```
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
```
<h5>blockchain/indexers/addrindex.go:</h5>

```
type AddrIndex struct{}
func NewAddrIndex(db database.DB, chainParams *chaincfg.Params) *AddrIndex {
	fmt.Println("Unfinished:indexers.NewAddrIndex")
	return nil
}
```

<h5>blockchain/indexers/cfindex.go:</h5>

```
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
```
<h5>blockchain/chain.go</h5>

```
// IndexManager 为索引器提供了一个通用接口
// provides a generic interface that the is called when blocks are
// connected and disconnected to and from the tip of the main chain for the
// purpose of supporting optional indexes.
type IndexManager interface{}
```
<h5>blockchain/indexers/manager.go</h5>

```
// Manager 定义一了个可以管理多个可选索引器的索引管理器,
// 同时实现了 blockchain.IndexManager 接口,这样它就能
// 无缝的传入普通区块链处理过程.
// defines an index manager that manages multiple optional indexes and
// implements the blockchain.IndexManager interface so it can be seamlessly
// plugged into normal chain processing.
type Manager struct {
	db             database.DB
	enabledIndexes []Indexer
}

// NewManager 返回一个索引管理器
// returns a new index manager with the provided indexes enabled.
//
// The manager returned satisfies the blockchain.IndexManager interface and thus
// cleanly plugs into the normal blockchain processing path.
func NewManager(db database.DB, enabledIndexes []Indexer) *Manager {
	return &Manager{
		db:             db,
		enabledIndexes: enabledIndexes,
	}
}
```
<h5>chaincfg/params.go</h5>

```
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
```
<h5>wire/protocol.go</h5>

```
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
```

<h4 id="204">区块链 blockchain</h4>
<h5>blockchain/chain.go</h5>

```
// Config is a descriptor which specifies the blockchain instance configuration.
type Config struct{
  // DB defines the database which houses the blocks and will be used to
	// store all metadata created by this package such as the utxo set.
	//
	// This field is required.
	DB database.DB

	// Interrupt specifies a channel the caller can close to signal that
	// long running operations, such as catching up indexes or performing
	// database migrations, should be interrupted.
	//
	// This field can be nil if the caller does not desire the behavior.
	Interrupt <-chan struct{}

	// ChainParams identifies which chain parameters the chain is associated
	// with.
	//
	// This field is required.
	ChainParams *chaincfg.Params

	// Checkpoints hold caller-defined checkpoints that should be added to
	// the default checkpoints in ChainParams.  Checkpoints must be sorted
	// by height.
	//
	// This field can be nil if the caller does not wish to specify any
	// checkpoints.
	Checkpoints []chaincfg.Checkpoint

	// TimeSource defines the median time source to use for things such as
	// block processing and determining whether or not the chain is current.
	//
	// The caller is expected to keep a reference to the time source as well
	// and add time samples from other peers on the network so the local
	// time is adjusted to be in agreement with other peers.
	TimeSource MedianTimeSource

	// SigCache defines a signature cache to use when when validating
	// signatures.  This is typically most useful when individual
	// transactions are already being validated prior to their inclusion in
	// a block such as what is usually done via a transaction memory pool.
	//
	// This field can be nil if the caller is not interested in using a
	// signature cache.
	SigCache *txscript.SigCache

	// IndexManager defines an index manager to use when initializing the
	// chain and connecting and disconnecting blocks.
	//
	// This field can be nil if the caller does not wish to make use of an
	// index manager.
	IndexManager IndexManager

	// HashCache defines a transaction hash mid-state cache to use when
	// validating transactions. This cache has the potential to greatly
	// speed up transaction validation as re-using the pre-calculated
	// mid-state eliminates the O(N^2) validation complexity due to the
	// SigHashAll flag.
	//
	// This field can be nil if the caller is not interested in using a
	// signature cache.
	HashCache *txscript.HashCache
}

// BlockChain 提供使用比特币区块链的功能
// provides functions for working with the bitcoin block chain.
// It includes functionality such as rejecting duplicate blocks, ensuring blocks
// follow all rules, orphan handling, checkpoint handling, and best chain
// selection with reorganization.
type BlockChain struct {
	stateLock     sync.RWMutex
	stateSnapshot *BestState
}
type BestState struct{
  Height      int32          // The height of the block.
}

// New returns a BlockChain instance using the provided configuration details.
func New(config *Config) (*BlockChain, error) {
	fmt.Println("Unfinished:blockchain.New")
	return nil, nil
}

// IndexManager 为索引器提供了一个通用接口
// provides a generic interface that the is called when blocks are
// connected and disconnected to and from the tip of the main chain for the
// purpose of supporting optional indexes.
type IndexManager interface{}

func (b *BlockChain) BestSnapshot() *BestState {
	b.stateLock.RLock()
	snapshot := b.stateSnapshot
	b.stateLock.RUnlock()
	return snapshot
}
```

<h4 id="205">数据库 database</h4>
<h5>database/interface.go</h5>

```
type DB interface {
	// Close 完全关闭数据库并同步所有数据.
	// 它会阻塞直到数据库所有事务都完成
	// cleanly shuts down the database and syncs all data.  It will
	// block until all database transactions have been finalized (rolled
	// back or committed).
	Close() error
	// Update invokes the passed function in the context of a managed
	// read-write transaction. 
	Update(fn func(tx Tx) error) error
}

// Tx 代表数据库事务.
// represents a database transaction.  It can either by read-only or
// read-write.  The transaction provides a metadata bucket against which all
// read and writes occur.
type Tx interface{}

```

<h4 id="206">交易池 mempool</h4>
<h5>mempool/mempool.go</h5>

```
type Config struct {
}

type TxPool struct{}

func New(cfg *Config) *TxPool {
	fmt.Println("Unfinished:mempool.New")
	return nil
}
```
<h5>mempool/estimatefee.go</h5>

```
type FeeEstimator struct {
	// The last known height
	lastKnownHeight int32

	mtx sync.RWMutex
}

// LastKnownHeight 返回最后注册的高度
// returns the height of the last block which was registered.
func (ef *FeeEstimator) LastKnownHeight() int32 {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	return ef.lastKnownHeight
}
```

<h4 id="207">同步管理 netsync</h4>
<h5>netsync/manager.go</h5>

```
type SyncManager struct{}
func New(config *Config) (*SyncManager, error) {
	fmt.Println("Unfinished:netsync.New")
	return nil, nil
}
```
<h5>netsync/interface.go</h5>

```
type Config struct{}
```

<h4 id="208">挖矿 mining</h4>
<h5>mining/policy.go</h5>

```
type Policy struct {
}
```
<h5>mining/mining.go</h5>

```
type BlkTmplGenerator struct {
}
type BlkTmplGenerator struct {
}
type TxSource interface {
}
func NewBlkTmplGenerator(policy *Policy, params *chaincfg.Params,
	txSource TxSource, chain *blockchain.BlockChain,
	timeSource blockchain.MedianTimeSource,
	sigCache *txscript.SigCache,
	hashCache *txscript.HashCache) *BlkTmplGenerator {
	fmt.Println("Unfinished:mining.BlkTmpGenerator")
	return nil
}
```
<h5>mining/mediantime.go</h5>

```
type MedianTimeSource interface{}
```
<h5>txscript/sigcache.go</h5>

```
type SigCache struct{}
```
<h5>txscript/hashcache.go</h5>

```
type HashCache struct{}
```
<h5>mining/cpuminer/cpuminer.go</h5>

```
type Config struct{
  // BlockTemplateGenerator 区块生成模板
	// identifies the instance to use in order to
	// generate block templates that the miner will attempt to solve.
	BlockTemplateGenerator *mining.BlkTmplGenerator
}
type CPUMiner struct{}
func New(cfg *Config) *CPUMiner{
	fmt.Println("Unfinished:cpuminer.New")
	return nil
}
```

<h4 id="209">连接管理 Connection Manager</h4>
<h5>connmgr/connmanager.go</h5>

```
type Config struct {
	// GetNewAddress 获取地址进行网络连接的方法.
	// is a way to get an address to make a network connection
	// to.  If nil, no new connections will be made automatically.
	GetNewAddress func() (net.Addr, error)
}
type ConnManager struct{}

// ConnReq 对某个网络连接地址的连接请求
// is the connection request to a network address. If permanent, the
// connection will be retried on disconnection.
type ConnReq struct {
	Addr net.Addr
	// Permanent 如果为true,当断开连接时会尝试重连
	Permanent bool
}

func New(cfg *Config) (*ConnManager, error) {
	fmt.Println("Unfinished:connmgr.New")
	return nil, nil
}
func (cm *ConnManager) Connect(c *ConnReq) {
	fmt.Println("Unfinished:connmgr.Connect")
}
```

<h4 id="210">RPC</h4>
<h5>rpcserver.go</h5>

```
type rpcserverConfig struct{}
type rpcServer struct {
	requestProcessShutdown chan struct{}
}

func newRPCServer(config *rpcserverConfig) (*rpcServer, error) {
	fmt.Println("Unfinished:newRPCServer")
	return nil, nil
}

// RequestedProcessShutdown 返回一个授信RPC客户端请求关闭进程时发送到的通道
// returns a channel that is sent to when an authorized
// RPC client requests the process to shutdown.  If the request can not be read
// immediately, it is dropped.
func (s *rpcServer) RequestedProcessShutdown() <-chan struct{} {
	return s.requestProcessShutdown
}
```