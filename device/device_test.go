package device

import (
	"context"
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Welasco/HubitatDeviceEvents/database"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/gofiber/fiber/v2"
)

type failingDatabase struct {
	database.Database
	err error
}

func (db failingDatabase) GetDevices(context.Context) ([]model.Device, error) {
	return nil, db.err
}
func (db failingDatabase) GetDeviceEvents(context.Context) ([]model.DeviceEvent, error) {
	return nil, db.err
}
func (db failingDatabase) RegisterDeviceEvent(ctx context.Context, _ *model.DeviceEvent) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return db.err
}

func TestHandlersDoNotAcknowledgeFailedOrCanceledOperations(t *testing.T) {
	original := db
	defer func() { db = original }()
	for _, tc := range []struct {
		name    string
		handler fiber.Handler
		err     error
		status  int
	}{
		{"devices read failure", GetDevices, errors.New("read failed"), 500},
		{"events read failure", GetDeviceEvents, errors.New("read failed"), 500},
		{"event deadline", RegisterDeviceEvent, context.DeadlineExceeded, 504},
		{"event cancellation", RegisterDeviceEvent, context.Canceled, 503},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db = failingDatabase{err: tc.err}
			app := fiber.New()
			app.Post("/", tc.handler)
			req := httptest.NewRequest("POST", "/", strings.NewReader(`{"content":{"deviceId":"42"}}`))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if resp.StatusCode != tc.status || !strings.Contains(string(body), `"status":"error"`) {
				t.Fatalf("failed operation returned %d: %s", resp.StatusCode, body)
			}
		})
	}
}
