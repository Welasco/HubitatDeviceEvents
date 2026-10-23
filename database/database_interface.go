package database

import (
	"context"

	"github.com/Welasco/HubitatDeviceEvents/model"
)

type Database interface {
	GetDevice(ctx context.Context, id string) (model.Device, error)

	GetDevices(ctx context.Context) ([]model.Device, error)

	AddDevice(ctx context.Context, device *model.Device) error

	DeleteDevice(ctx context.Context, id string) error

	UpdateDevice(ctx context.Context, device *model.Device) error

	GetDeviceEvents(ctx context.Context) ([]model.DeviceEvent, error)

	GetDeviceEventId(ctx context.Context, id string, timeStamp model.UrlQueries) ([]model.DeviceEvent, error)

	RegisterDeviceEvent(ctx context.Context, device *model.DeviceEvent) error

	Close() error
}
