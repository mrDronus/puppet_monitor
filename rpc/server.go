package rpc

import "puppet_monitoring/impl"

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"syscall"
	"time"
)

type RPCServer struct {
	Envs *impl.EnvironmentCollection
}

var server *RPCServer

func (s *RPCServer) CreateServer(conf impl.Settings) {
	server = s
	ns := new(PPTRpc)
	rpc.Register(ns)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", conf.Ip+":"+strconv.Itoa(conf.RpcPort))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func (t *PPTRpc) RemoveNode(args *RemoveNodeArgs, reply *string) error {
	log.Println("[REMOTE COMMAND] RemoveNode: " + args.Host)
	*reply = server.Envs.RemoveNode(args.Host)
	return nil
}

func (t *PPTRpc) GetStatus(args *GetStatusArgs, reply *string) error {
	log.Println("[REMOTE COMMAND] GetStatus")
	*reply = server.Envs.ProcessCollectionState(args.Errors)
	return nil
}

func (t *PPTRpc) StopMasterProcess(args *EmptyArgs, reply *bool) error {
	log.Println("[REMOTE COMMAND] Stop master process")
	*reply = true
	go func() {
		time.Sleep(1 * time.Second)
		current_pid := os.Getpid()
		proc, _ := os.FindProcess(current_pid)
		log.Printf("sending TERM signal to pid: %v\n", current_pid)
		proc.Signal(syscall.SIGTERM)
	}()
	return nil
}

func (t *PPTRpc) GetInfo(args *EmptyArgs, reply *string) error {
	log.Println("[REMOTE COMMAND] GetInfo")
	*reply = server.Envs.GetInfo()
	return nil
}
