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
