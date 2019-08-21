package addrmgr

import (
	"net"

	"github.com/iblockchains/bitcoin/wire"
)

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