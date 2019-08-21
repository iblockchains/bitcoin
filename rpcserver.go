package main

import "fmt"

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
