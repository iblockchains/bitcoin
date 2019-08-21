package addrmgr

import (
	"fmt"
	"net"
	"testing"

	"github.com/iblockchains/bitcoin/wire"
)

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
