package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/emadolsky/automaxprocs/maxprocs"
	"github.com/tchorzewski1991/fds/base/logger"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/zap"
)

var build = "develop"
var service = "FLIGHTS-API"

func main() {
	// ================================================================================================================
	// Set GOMAXPROCS

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	if _, err := maxprocs.Set(); err != nil {
		fmt.Println("maxprocs: %w", err)
	}
	cpu := runtime.GOMAXPROCS(0)

	// ================================================================================================================
	// Construct application logger

	// Set logger fields common to all logs
	fields := []logger.Field{
		{
			Name: "service", Value: service,
		},
		{
			Name: "build", Value: build,
		},
		{
			Name: "cpu", Value: cpu,
		},
	}
	// Create new logger
	l, err := logger.New(fields...)
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	// Flush internal logger buffer
	defer func() {
		if err = l.Sync(); err != nil {
			fmt.Println("Error flushing logger:", err)
		}
	}()
	// Run application
	if err = run(l); err != nil {
		l.Errorf("Error running app: %s", err)
		os.Exit(1)
	}
}

func run(logger *zap.SugaredLogger) error {
	logger.Info("Starting service")
	defer logger.Info("Stopped service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	logger.Info("Stopping service")

	return nil
}
