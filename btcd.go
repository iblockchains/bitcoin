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
