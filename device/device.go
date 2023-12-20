package device

import (
	"strconv"

	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/gofiber/fiber/v2"
)

var db database.Database

func DBInit(config model.Config) {
	// Check if database is already initialized
	if config.DatabaseType == "mysql" {
		db, err := database.NewMysqlDB("ConnectionString", "DatabaseName", 1433, "UserName", "Password")
		if err != nil {
			panic(err)
		}
		db = db
	}
	// Implement additional database types here
}

func GetDevices(c *fiber.Ctx) error {
	//var devices []model.Device
	devices, _ := db.GetDevices()
	return c.JSON(devices)
}

func GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	int_id, _ := strconv.Atoi(id)
	var device model.Device
	//db.Find(&device, id)
	device, _ = db.GetDevice(int_id)
	return c.JSON(device)
}

func AddDevice(c *fiber.Ctx) error {
	var device model.Device
	if err := c.BodyParser(device); err != nil {
		//c.Status(503).Send(err)
		c.Status(503).SendString(err.Error())
		return err
	}

	db.AddDevice(&device)
	return c.JSON(device)
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	int_id, _ := strconv.Atoi(id)
	var device model.Device
	device, _ = db.GetDevice(int_id)
	if device.DisplayName == "" {
		c.Status(500).SendString("No device found with given ID")
		return nil
	}
	db.DeleteDevice(int_id)
	return c.SendString("Device successfully deleted")
}

func UpdateDevice(c *fiber.Ctx) error {
	var device model.Device
	if err := c.BodyParser(device); err != nil {
		//c.Status(503).Send(err)
		c.Status(503).SendString(err.Error())
		return err
	}

	db.UpdateDevice(&device)
	return c.JSON(device)
}
