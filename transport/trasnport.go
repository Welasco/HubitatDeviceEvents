package transport

import (
	"github.com/Welasco/HubitatDeviceEvents/device"
	//"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/v2"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/api/v1/device", device.GetDevices)
	app.Get("/api/v1/device/:isin", device.GetDevice)
	app.Post("/api/v1/device", device.NewDevice)
	app.Delete("/api/v1/device/:isin", device.DeleteDevice)
}

// Setup - set's up our fiber app and the routes
// returns a pointer to app
func Setup() *fiber.App {
	app := fiber.New()
	setupRoutes(app)
	return app
}
