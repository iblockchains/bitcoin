package connmgr

import (
	"fmt"
	"net"
)

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
