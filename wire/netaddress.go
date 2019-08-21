package wire

import (
	"net"
	"time"
)

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
