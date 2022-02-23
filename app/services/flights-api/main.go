package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/emadolsky/automaxprocs/maxprocs"
	_ "go.uber.org/automaxprocs"
)

var build = "develop"
var service = "flights-api"

func main() {
	// =========================================================================
	// GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
		os.Exit(1)
	}

	cpu := runtime.GOMAXPROCS(0)

	log.Printf("Starting service[%s] build[%s] cpu[%d]\n", service, build, cpu)
	defer log.Printf("Stopped service[%s] build[%s] cps[%d]\n", service, build, cpu)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	log.Printf("Stopping service[%s] build[%s] cpu[%d]\n", service, build, cpu)
}
