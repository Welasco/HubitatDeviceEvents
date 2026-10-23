package transport

import (
	"context"
	"os"
	"time"

	"github.com/Welasco/HubitatDeviceEvents/device"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func setupRoutes(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/api/v1/device/event", device.GetDeviceEvents)
	app.Post("/api/v1/device/event", device.RegisterDeviceEvent)
	app.Get("/api/v1/device/:id/event", device.GetDeviceEventId)
	app.Get("/api/v1/device", device.GetDevices)
	app.Get("/api/v1/device/:id", device.GetDevice)
	app.Post("/api/v1/device", device.AddDevice)
	app.Delete("/api/v1/device/:id", device.DeleteDevice)
	app.Put("/api/v1/device", device.UpdateDevice)
}

// Setup - set's up our fiber app and the routes
// returns a pointer to app
func Setup(requestContext context.Context) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	})
	app.Use(func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(requestContext, 5*time.Second)
		defer cancel()
		c.SetUserContext(ctx)
		return c.Next()
	})

	if os.Getenv("Environment") == "Dev" {
		app.Static("/", "./frontend/dist")
	} else {
		app.Static("/", "./public")
	}
	app.Use(cors.New())
	setupRoutes(app)

	return app
}

// _ "github.com/joho/godotenv/autoload"
// )

// var db database.Database

// var config model.Config

// func LoadConfig() {
// 	config = model.Config{
// 		ConnectionString: os.Getenv("ConnectionString"),
// 		DatabaseType:     os.Getenv("DatabaseType"),
// 	}
// }
