package main

import (
	"fmt"
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
	app := transport.Setup()
	device.DBInit()
	fmt.Println("Hello World")
	logger.Info("Starting server on port 3000")
	app.Listen(":3000")
}
