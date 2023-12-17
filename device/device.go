package device

import (
	"fmt"

	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/gofiber/fiber/v2"
	//"github.com/jinzhu/gorm"
)

type Device struct {
	//gorm.Model
	ISIN   int    `json:"ISIN"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Rating int    `json:"rating"`
}

func GetDevices(c *fiber.Ctx) {
	db := database.DBConn
	var devices []Device
	db.Find(&devices)
	c.JSON(devices)
}

func GetDevice(c *fiber.Ctx) {
	id := c.Params("id")
	db := database.DBConn
	var device Device
	db.Find(&device, id)
	c.JSON(device)
}

func NewDevice(c *fiber.Ctx) {
	db := database.DBConn

	device := new(Device)
	if err := c.BodyParser(device); err != nil {
		c.Status(503).Send(err)
		return
	}

	db.Create(&device)
	c.JSON(device)
}

func DeleteDevice(c *fiber.Ctx) {
	ISIN := c.Params("isin")
	db := database.DBConn

	fmt.Println(ISIN)

	var device Device
	db.Find(&device, "ISIN = ?", ISIN)
	if device.Title == "" {
		c.Status(500).Send("No device found with given ID")
		return
	}
	db.Delete(&device)
	c.Send("Device successfully deleted")
}
