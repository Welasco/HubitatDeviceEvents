package device

import (
	"os"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

var db database.Database

var config model.Config

func LoadConfig() {
	config = model.Config{
		ConnectionString: os.Getenv("ConnectionString"),
		DatabaseType:     os.Getenv("DatabaseType"),
	}
}

func DBInit() {
	LoadConfig()
	logger.Info("[device][DBInit] Initializing database")
	if config.DatabaseType == "mysql" {
		var err error
		db, err = database.NewMySQLDB(config.ConnectionString)
		if err != nil {
			logger.Error("[device][DBInit] Error initializing database")
			logger.Error("[device][DBInit] " + err.Error())
			//panic(err)
		}
	} else {
		logger.Error("[device][DBInit] Database type " + config.DatabaseType + " not supported")
	}

	// Implement additional database types here
}

func GetDevices(c *fiber.Ctx) error {
	devices, err := db.GetDevices()
	if err != nil {
		logger.Error("[device][GetDevices] Error reading devices from database")
		logger.Error("[device][GetDevices] Error: " + err.Error())
	}

	return c.JSON(devices)
}

func GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device model.Device
	device, err := db.GetDevice(id)
	if err != nil {
		logger.Error("[device][GetDevice] Error reading devices from database")
		logger.Error("[device][GetDevice] Error: " + err.Error())
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Device not found", "data": err.Error()})
	}
	return c.JSON(device)
}

func AddDevice(c *fiber.Ctx) error {
	var device model.Device
	if err := c.BodyParser(&device); err != nil {
		logger.Error("[device][AddDevice] Error adding device")
		logger.Error("[device][AddDevice] Error: " + err.Error())
		c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing body", "data": err.Error()})
		return err
	}
	err := db.AddDevice(&device)
	if err != nil {
		logger.Error("[device][AddDevice] Error adding device")
		logger.Error("[device][AddDevice] Error: " + err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error adding device", "data": err.Error()})
	}
	return c.JSON(device)
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device model.Device
	device, err := db.GetDevice(id)
	if err != nil {
		logger.Error("[device][DeleteDevice] Error checking device existance before deleting")
		logger.Error("[device][DeleteDevice] Error: " + err.Error())
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Device not found", "data": err.Error()})
	}
	if device.Id == "" {
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Device not found", "data": err.Error()})
	}
	err = db.DeleteDevice(id)
	if err != nil {
		logger.Error("[device][DeleteDevice] Error deleting device")
		logger.Error("[device][DeleteDevice] Error: " + err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unable to delete device, check if there are no events associated", "data": err.Error()})
	}
	return c.JSON(fiber.Map{"status": "success", "message": "Device deleted"})
}

func UpdateDevice(c *fiber.Ctx) error {
	var device model.Device
	if err := c.BodyParser(&device); err != nil {
		logger.Error("[device][UpdateDevice] Error updating device")
		logger.Error("[device][UpdateDevice] Error: " + err.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing body", "data": err.Error()})
	}

	err := db.UpdateDevice(&device)
	if err != nil {
		logger.Error("[device][UpdateDevice] Error deleting device")
		logger.Error("[device][UpdateDevice] Error: " + err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Unable to delete device, check if there are no events associated", "data": err.Error()})
	}
	return c.JSON(device)
}

func GetDeviceEvents(c *fiber.Ctx) error {
	deviceevents, err := db.GetDeviceEvents()
	if err != nil {
		logger.Error("[device][GetDeviceEvents] Error reading devices from database")
		logger.Error("[device][GetDeviceEvents] Error: " + err.Error())
	}

	return c.JSON(deviceevents)
}

func GetDeviceEventId(c *fiber.Ctx) error {
	id := c.Params("id")
	var deviceevents []model.DeviceEvent
	deviceevents, err := db.GetDeviceEventId(id)
	if err != nil {
		logger.Error("[device][GetDeviceEventId] Error reading events from device from database")
		logger.Error("[device][GetDeviceEventId] Error: " + err.Error())
		return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Device not found", "data": err.Error()})
	}
	return c.JSON(deviceevents)
}

func RegisterDeviceEvent(c *fiber.Ctx) error {
	var devicecontent model.Content
	logger.Debug("[device][RegisterDeviceEvent] Received device events")
	logger.Debug("[device][RegisterDeviceEvent] Event: " + string(c.BodyRaw()))
	if err := c.BodyParser(&devicecontent); err != nil {
		logger.Error("[device][RegisterDeviceEvent] Error reading devices from database")
		logger.Error("[device][RegisterDeviceEvent] Error: " + err.Error())
		return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Error parsing body", "data": err.Error()})
	}
	err := db.RegisterDeviceEvent(&devicecontent.DeviceEvent)
	if err != nil {
		logger.Error("[device][RegisterDeviceEvent] Error Registering device event")
		logger.Error("[device][RegisterDeviceEvent] Error: " + err.Error())
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Error registering device event", "data": err.Error()})
	}
	return c.SendString("{}")
}
