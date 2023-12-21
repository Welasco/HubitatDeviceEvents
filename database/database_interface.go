package database

import (
	//"database/sql"

	"github.com/Welasco/HubitatDeviceEvents/model"
)

type Database interface {
	GetDevice(id int) (model.Device, error)

	GetDevices() ([]model.Device, error)

	AddDevice(device *model.Device) error

	DeleteDevice(id int) error

	UpdateDevice(device *model.Device) error

	GetDeviceEvents() error
}
