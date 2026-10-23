package device

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

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

func DBInit(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	LoadConfig()
	logger.Info("[device][DBInit] Initializing database")
	if config.DatabaseType == "mysql" {
		mysqlDB, err := database.NewMySQLDB(ctx, config.ConnectionString)
		if err != nil {
			logger.Error("[device][DBInit] Error initializing database")
			logger.Error("[device][DBInit] " + err.Error())
			return err
		}
		db = mysqlDB
	} else {
		logger.Error("[device][DBInit] Database type " + config.DatabaseType + " not supported")
		return fmt.Errorf("database type %q not supported", config.DatabaseType)
	}

	return nil
}

func DBClose() error {
	if db == nil {
		return nil
	}
	logger.Info("[device][DBClose] Closing database connection")
	if err := db.Close(); err != nil {
		logger.Error("[device][DBClose] Error closing database connection")
		logger.Error("[device][DBClose] " + err.Error())
		return err
	}
	return nil
}

func databaseError(c *fiber.Ctx, status int, message string, err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		status = fiber.StatusGatewayTimeout
	} else if errors.Is(err, context.Canceled) {
		status = fiber.StatusServiceUnavailable
	}
	return c.Status(status).JSON(fiber.Map{"status": "error", "message": message, "data": err.Error()})
}

func GetDevices(c *fiber.Ctx) error {
	devices, err := db.GetDevices(c.UserContext())
	if err != nil {
		logger.Error("[device][GetDevices] Error reading devices from database")
		logger.Error("[device][GetDevices] Error: " + err.Error())
		return databaseError(c, 500, "Error reading devices", err)
	}

	return c.JSON(devices)
}

func GetDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device model.Device
	device, err := db.GetDevice(c.UserContext(), id)
	if err != nil {
		logger.Error("[device][GetDevice] Error reading devices from database")
		logger.Error("[device][GetDevice] Error: " + err.Error())
		return databaseError(c, 404, "Device not found", err)
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
	err := db.AddDevice(c.UserContext(), &device)
	if err != nil {
		logger.Error("[device][AddDevice] Error adding device")
		logger.Error("[device][AddDevice] Error: " + err.Error())
		return databaseError(c, 500, "Error adding device", err)
	}
	return c.JSON(device)
}

func DeleteDevice(c *fiber.Ctx) error {
	id := c.Params("id")
	var device model.Device
	device, err := db.GetDevice(c.UserContext(), id)
	if err != nil {
		logger.Error("[device][DeleteDevice] Error checking device existance before deleting")
		logger.Error("[device][DeleteDevice] Error: " + err.Error())
		return databaseError(c, 404, "Device not found", err)
	}
	if device.Id == "" {
		return databaseError(c, 404, "Device not found", sql.ErrNoRows)
	}
	err = db.DeleteDevice(c.UserContext(), id)
	if err != nil {
		logger.Error("[device][DeleteDevice] Error deleting device")
		logger.Error("[device][DeleteDevice] Error: " + err.Error())
		return databaseError(c, 500, "Unable to delete device, check if there are no events associated", err)
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

	err := db.UpdateDevice(c.UserContext(), &device)
	if err != nil {
		logger.Error("[device][UpdateDevice] Error deleting device")
		logger.Error("[device][UpdateDevice] Error: " + err.Error())
		return databaseError(c, 500, "Unable to delete device, check if there are no events associated", err)
	}
	return c.JSON(device)
}

func GetDeviceEvents(c *fiber.Ctx) error {
	deviceevents, err := db.GetDeviceEvents(c.UserContext())
	if err != nil {
		logger.Error("[device][GetDeviceEvents] Error reading devices from database")
		logger.Error("[device][GetDeviceEvents] Error: " + err.Error())
		return databaseError(c, 500, "Error reading device events", err)
	}

	return c.JSON(deviceevents)
}

func GetDeviceEventId(c *fiber.Ctx) error {
	id := c.Params("id")
	var deviceevents []model.DeviceEvent
	var err error
	urlQueries := c.Queries()
	var urlQueriesStuct model.UrlQueries
	urlQueriesStuct.Start = urlQueries["start"]
	urlQueriesStuct.End = urlQueries["end"]

	deviceevents, err = db.GetDeviceEventId(c.UserContext(), id, urlQueriesStuct)
	if err != nil {
		logger.Error("[device][GetDeviceEventId] Error reading events from device from database")
		logger.Error("[device][GetDeviceEventId] Error: " + err.Error())
		return databaseError(c, 404, "Device not found or invalid query string start or end date.", err)
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
	err := db.RegisterDeviceEvent(c.UserContext(), &devicecontent.DeviceEvent)
	if err != nil {
		logger.Error("[device][RegisterDeviceEvent] Error Registering device event")
		logger.Error("[device][RegisterDeviceEvent] Error: " + err.Error())
		return databaseError(c, 500, "Error registering device event", err)
	}
	return c.SendString("{}")
}
