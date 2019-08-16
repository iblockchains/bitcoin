<h1>搭建基本的启动服务框架 start a simple server</h1>

* [第一步](#2)

  * [主函数 btcd](#2.1)

  * [日志 log](#2.2)

  * [配置 config](#2.3)

  * [信号 signal](#2.4)

  * [升级 upgrade](#2.5)

  * [服务器 server](#2.6)

  * [版本信息 version](#2.7)

  * [参数](#2.8)

* [第二步](#3)
  * [修改服务 update server](#3.1)

<h4 id="2.1">主函数 main</h4>
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

<h4 id="2.2">日志 log</h4>
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

<h4 id="2.3">配置 config</h4>
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
	Listeners      []string `long:"listen" description:"Add an interface/port to listen for connections (default all interfaces port: 8333, testnet: 18333)"`
	AgentBlacklist []string `long:"agentblacklist" description:"A comma separated list of user-agent substrings which will cause btcd to reject any peers whose user-agent contains any of the blacklisted substrings."`
	AgentWhitelist []string `long:"agentwhitelist" description:"A comma separated list of user-agent substrings which will cause btcd to require all peers' user-agents to contain one of the whitelisted substrings. The blacklist is applied before the blacklist, and an empty whitelist will allow all agents that do not fail the blacklist."`
	LogDir         string   `long:"logdir" description:"Directory to log output."`
	Profile        string   `long:"profile" description:"Enable HTTP profiling on given port -- NOTE port must be between 1024 and 65536"`
	CPUProfile     string   `long:"cpuprofile" description:"Write CPU profile to the specified file"`
	DropAddrIndex  bool     `long:"dropaddrindex" description:"Deletes the address-based transaction index from the database on start up and then exits."`
	DropTxIndex    bool     `long:"droptxindex" description:"Deletes the hash-based transaction index from the database on start up and then exits."`
	DropCfIndex    bool     `long:"dropcfindex" description:"Deletes the index used for committed filtering (CF) support from the database on start up and then exits."`
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

<h4 id="2.4">信号 signal</h4>
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

<h4 id="2.5">升级 upgrade</h4>
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

<h4 id="2.6">服务器 server</h4>
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

<h4 id="2.7">版本信息</h4>
<h5>新建 version.go</h5>

```
package main

import "fmt"

func version() string {
	fmt.Println("Unfinished:version")
	return "version"
}
```

<h4 id="2.8">app参数</h4>
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

<h4 id="3.1">修改服务器 update server</h4>