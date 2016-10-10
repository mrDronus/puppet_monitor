package main

import (
	"log"
	"os"
)

const (
	VERSION = "1.0.0 \u00a9Inventos, Orel (RU), 2016"
)

func main() {
	log.SetOutput(os.Stdout)
	if process_args() {
		os.Exit(0)
	}
	run_master_process()
}
