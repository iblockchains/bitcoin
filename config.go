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
	fmt.Println("待:loadConfig")
	// Initialize log rotation.  After log rotation has been initialized, the
	// logger variables may be used.
	initLogRotator(filepath.Join(cfg.LogDir, defaultLogFilename))
	return &cfg, nil, nil
}
