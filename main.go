package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/device"
	"github.com/Welasco/HubitatDeviceEvents/transport"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	file := os.Getenv("logPath") + os.Getenv("HOSTNAME") + ".log"
	logger.Init(file, os.Getenv("logLevel"))
}

func main() {
	logger.Info("[main][main] Starting Hubitat Device Events")
	app := transport.Setup()
	device.DBInit()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		logger.Info("[main][main] Starting server on port " + port)
		if err := app.Listen(":" + port); err != nil {
			logger.Error("[main][main] Server stopped with error: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Info("[main][main] Received signal: " + sig.String() + ", shutting down gracefully")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("[main][main] Error during server shutdown: " + err.Error())
	} else {
		logger.Info("[main][main] Server stopped accepting new connections")
	}

	device.DBClose()

	logger.Info("[main][main] Shutdown complete")
}
