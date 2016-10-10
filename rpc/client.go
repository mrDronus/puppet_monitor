package rpc

import (
	"log"
	"net/rpc"
	"os"
	"puppet_monitoring/impl"
)

type RPCClient struct {
	Conf *impl.Settings
}

func (p *RPCClient) GetStatus(with_errors bool) string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.GetStatus", &GetStatusArgs{Errors: with_errors}, &result)
	return result
}

func (p *RPCClient) RemoveNode(host string) string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.RemoveNode", &RemoveNodeArgs{Host: host}, &result)
	return result
}

func (p *RPCClient) GetInfo() string {
	cl := p.createClient()
	var result string
	cl.Call("PPTRpc.GetInfo", &EmptyArgs{}, &result)
	return result
}

func (p *RPCClient) StopMasterProcess() (bool, error) {
	cl := p.createClient()
	var result bool
	err := cl.Call("PPTRpc.StopMasterProcess", &EmptyArgs{}, &result)
	return result && err == nil, err
}

func (c RPCClient) createClient() *rpc.Client {
	client, err := rpc.DialHTTP("tcp", c.Conf.RpcComputed)
	if err != nil {
		log.Fatal("connect to master process failed:", err)
		os.Exit(1)
	}
	return client
}
