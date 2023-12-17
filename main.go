package main

import (
	"fmt"

	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/Welasco/HubitatDeviceEvents/transport"
)

func main() {
	app := transport.Setup()
	database.InitDatabase()
	fmt.Println("Hello World")
	app.Listen(3000)
}
