package main

import (
	"os"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/device"
	"github.com/Welasco/HubitatDeviceEvents/transport"
	_ "github.com/joho/godotenv/autoload"
)

func init() {
	logger.Init(os.Getenv("logPath"), os.Getenv("logLevel"))
}

func main() {
	logger.Info("[main][main] Starting Hubitat Device Events")
	app := transport.Setup()
	device.DBInit()
	logger.Info("[main][main] Starting server on port 3000")
	app.Listen(":3000")
}
