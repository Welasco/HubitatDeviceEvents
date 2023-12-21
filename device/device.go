package device

import (
	"strconv"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/gofiber/fiber/v2"
)

var db database.Database

func DBInit(config model.Config) {
	// Check if database is already initialized
	if config.DatabaseType == "mysql" {
		var err error
		db, err = database.NewMysqlDB(config.ConnectionString, config.DatabaseName, 1433, config.UserName, config.Password)
		if err != nil {
			logger.Error("[device][DBInit] Error initializing database")
			logger.Error("[device][DBInit] " + err.Error())
			//panic(err)
		}
	}
	// Implement additional database types here
}

func GetDevices(c *fiber.Ctx) error {
	//var devices []model.Device
	devices, err := db.GetDevices()
	if err != nil {
		logger.Info("[device][GetDevices] Error reading devices from database")
		logger.Info("[device][GetDevices] Error: " + err.Error())
	}

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
	if err := c.BodyParser(&device); err != nil {
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
	if err := c.BodyParser(&device); err != nil {
		//c.Status(503).Send(err)
		c.Status(503).SendString(err.Error())
		return err
	}

	db.UpdateDevice(&device)
	return c.JSON(device)
}
