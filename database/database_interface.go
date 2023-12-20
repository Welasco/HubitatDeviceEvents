package database

import (
	"database/sql"

	"github.com/Welasco/HubitatDeviceEvents/model"
)

type Database interface {
	InitDB() (*sql.DB, error)

	GetDevice(id int) (model.Device, error)

	GetDevices() ([]model.Device, error)

	AddDevice(device *model.Device) string

	DeleteDevice(id int) string

	UpdateDevice(device *model.Device) string

	GetDeviceEvents() error
}
