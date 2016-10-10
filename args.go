package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"puppet_monitoring/rpc"
)

var opts struct {
	Print   bool   `short:"p" long:"print" description:"print current environment collections"`
	Status  bool   `short:"s" long:"status" description:"print status"`
	Errors  bool   `short:"e" long:"error" description:"print status with errors"`
	Version bool   `short:"v" long:"version" description:"print version"`
	Stop    bool   `long:"stop" description:"send selfkill signal to master process"`
	Remove  string `short:"r" long:"remove" description:"remove all data about specified host"`
	Rpc     string `long:"rpc" description:"set rpc params for master process communication"`
}

func parse_args() {
	flags.Parse(&opts)
}

func process_args() bool {

	parse_args()

	if opts.Rpc != "" {
		settings.RpcComputed = opts.Rpc
	}

	if opts.Print {
		client := rpc.RPCClient{Conf: &settings}
		fmt.Println(client.GetInfo())
		os.Exit(0)
	}

	if opts.Status {
		client := rpc.RPCClient{Conf: &settings}
		fmt.Println(client.GetStatus(opts.Errors))
		os.Exit(0)
	}

	if opts.Remove != "" {
		client := rpc.RPCClient{Conf: &settings}
		fmt.Println(client.RemoveNode(opts.Remove))
		os.Exit(0)
	}

	if opts.Version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	if opts.Stop {
		client := rpc.RPCClient{Conf: &settings}
		result, err := client.StopMasterProcess()
		if result {
			fmt.Println("OK")
			os.Exit(0)
		} else {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	return len(os.Args[1:]) > 0
}
