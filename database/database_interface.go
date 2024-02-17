package database

import (
	"github.com/Welasco/HubitatDeviceEvents/model"
)

type Database interface {
	GetDevice(id string) (model.Device, error)

	GetDevices() ([]model.Device, error)

	AddDevice(device *model.Device) error

	DeleteDevice(id string) error

	UpdateDevice(device *model.Device) error

	GetDeviceEvents() ([]model.DeviceEvent, error)

	GetDeviceEventId(id string, timeStamp model.UrlQueries) ([]model.DeviceEvent, error)

	RegisterDeviceEvent(device *model.DeviceEvent) error
}
