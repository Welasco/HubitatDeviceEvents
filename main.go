package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/device"
	"github.com/Welasco/HubitatDeviceEvents/transport"
	_ "github.com/joho/godotenv/autoload"
)

const (
	drainTimeout    = 10 * time.Second
	shutdownTimeout = 15 * time.Second
)

func main() {
	os.Exit(run())
}

func run() int {
	var file string
	if value := os.Getenv("logToFile"); value != "" {
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			logger.Error("[main][run] Invalid logToFile: " + err.Error())
			return 1
		}
		if enabled {
			file = os.Getenv("logPath") + os.Getenv("HOSTNAME") + ".log"
		}
	}
	if err := logger.Init(file, os.Getenv("logLevel")); err != nil {
		logger.Error("[main][run] Error initializing logger: " + err.Error())
		return 1
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(quit)
	requestContext, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()

	logger.Info("[main][main] Starting Hubitat Device Events")
	app := transport.Setup(requestContext)
	initialized := make(chan error, 1)
	go func() {
		initialized <- device.DBInit(requestContext)
	}()
	select {
	case sig := <-quit:
		cancelRequests()
		return finishShutdown("Received signal during startup: "+sig.String(), func(ctx context.Context) error {
			select {
			case err := <-initialized:
				if errors.Is(err, context.Canceled) {
					return nil
				}
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		}, cancelRequests)
	case err := <-initialized:
		if err != nil {
			finishShutdown("Database initialization failed", func(context.Context) error { return err }, cancelRequests)
			return 1
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	listener, err := net.Listen("tcp4", ":"+port)
	if err != nil {
		finishShutdown("Unable to listen", func(context.Context) error { return err }, cancelRequests)
		return 1
	}
	serverResult := make(chan error, 1)
	go func() {
		serverResult <- app.Listener(listener)
	}()

	drain := func(ctx context.Context) error {
		err := app.ShutdownWithContext(ctx)
		// Also close our listener if shutdown raced with the start of app.Listener.
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			err = errors.Join(err, closeErr)
		}
		return err
	}
	select {
	case sig := <-quit:
		return finishShutdown("Received signal: "+sig.String(), drain, cancelRequests)
	case err := <-serverResult:
		finishShutdown("Server stopped unexpectedly", func(ctx context.Context) error {
			return errors.Join(fmt.Errorf("server stopped unexpectedly: %v", err), drain(ctx))
		}, cancelRequests)
		return 1
	}
}

func finishShutdown(reason string, drain func(context.Context) error, cancelRequests context.CancelFunc) int {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := shutdown(ctx, func(ctx context.Context) error {
		logger.Info("[main][shutdown] " + reason + ", shutting down gracefully")
		return drain(ctx)
	}, cancelRequests, device.DBClose)
	if err != nil {
		// Logging must not extend the exit deadline if stdout or optional file I/O is stuck.
		reported := make(chan struct{})
		go func() {
			logger.Error("[main][shutdown] Shutdown failed: " + err.Error())
			close(reported)
		}()
		select {
		case <-reported:
		case <-ctx.Done():
		}
		return 1
	}
	return 0
}

func shutdown(ctx context.Context, drain func(context.Context) error, cancelRequests context.CancelFunc, closeDB func() error) error {
	defer cancelRequests()
	finished := make(chan error, 1)
	go func() {
		start := time.Now()
		drainContext, cancel := context.WithTimeout(ctx, drainTimeout)
		err := drain(drainContext)
		cancel()
		cancelRequests()
		if err != nil {
			logger.Error("[main][shutdown] Drain failed: " + err.Error())
		}
		logger.Info(fmt.Sprintf("[main][shutdown] HTTP drain finished after %s; closing database", time.Since(start)))
		closeErr := closeDB()
		if closeErr != nil {
			logger.Error("[main][shutdown] Database close failed: " + closeErr.Error())
		}
		err = errors.Join(err, closeErr, ctx.Err())
		if err == nil {
			logger.Info(fmt.Sprintf("[main][shutdown] Shutdown complete after %s", time.Since(start)))
		}
		finished <- err
	}()

	select {
	case err := <-finished:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
