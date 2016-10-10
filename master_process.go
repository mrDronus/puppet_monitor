package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
)

import i "puppet_monitoring/impl"
import rpc "puppet_monitoring/rpc"

var settings = i.Settings{}.LoadSettings()

func run_master_process() {

	log.Printf("PID:%v\n", os.Getpid())
	runtime.GOMAXPROCS(1)
	if check_pid_file() {
		fmt.Println(settings.PidFile + " already exists! (other instance run?)")
		os.Exit(1)
	}

	create_pid_file()
	defer kill_pid()

	var laddr, err = net.ResolveTCPAddr("tcp", settings.Ip+":"+strconv.Itoa(settings.Port))
	ln, err := net.ListenTCP("tcp", laddr)

	if err != nil {
		panic(err)
	}

	defer ln.Close()
	log.Println("listening on", ln.Addr())
	service := i.Service{}.NewService()
	envs := i.EnvironmentCollection{}.NewEnvironmentCollection()
	envs.Conf = &settings
	service.SetEnvCollection(&envs)
	go service.HandleListener(ln)

	rpcsrv := rpc.RPCServer{Envs: &envs}
	rpcsrv.CreateServer(settings)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	var sig os.Signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig = <-ch:
		log.Println(sig)
	}

	service.Stop()
}

func check_pid_file() bool {
	var _, err = os.Stat(settings.PidFile)
	return err == nil
}

func create_pid_file() {
	var fd, err = os.Create(settings.PidFile)
	fd.Close()
	if err != nil {
		fmt.Println("Error creating pid file!")
		panic(err)
	}
	os.Chmod(settings.PidFile, 0644)
	fd, err = os.OpenFile(settings.PidFile, os.O_RDWR, 0644)
	defer fd.Close()
	var _, werr = fd.WriteString(strconv.Itoa(os.Getpid()))
	if werr != nil {
		fmt.Println("Error write pid file!")
		panic(werr)
	}
	fd.Sync()
}

func kill_pid() {
	err := os.Remove(settings.PidFile)
	if err != nil {
		log.Println(err)
	}
}
