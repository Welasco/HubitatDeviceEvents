package main

import (
	"fmt"
	"os"

	"github.com/Welasco/HubitatDeviceEvents/device"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/Welasco/HubitatDeviceEvents/transport"
	_ "github.com/joho/godotenv/autoload"
)

var config model.Config

func LoadConfig() {
	config = model.Config{
		ConnectionString: os.Getenv("ConnectionString"),
		DatabaseName:     os.Getenv("DatabaseName"),
		Port:             os.Getenv("Port"),
		UserName:         os.Getenv("UserName"),
		Password:         os.Getenv("Password"),
		DatabaseType:     os.Getenv("DatabaseType"),
	}
}

func init() {
	LoadConfig()

}

func main() {
	app := transport.Setup()
	device.DBInit(config)
	fmt.Println("Hello World")
	app.Listen("3000")
}
