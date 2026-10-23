package transport

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TestRequestContextHasDeadlineAndReceivesShutdownCancellation(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	app := Setup(parent)
	entered := make(chan context.Context, 1)
	app.Get("/test-context", func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		entered <- ctx
		<-ctx.Done()
		return c.SendStatus(fiber.StatusServiceUnavailable)
	})
	response := make(chan error, 1)
	go func() {
		resp, err := app.Test(httptest.NewRequest("GET", "/test-context", nil))
		if err == nil {
			defer resp.Body.Close()
			if resp.StatusCode != fiber.StatusServiceUnavailable {
				t.Errorf("unexpected response: %d", resp.StatusCode)
			}
		}
		response <- err
	}()
	select {
	case ctx := <-entered:
		deadline, ok := ctx.Deadline()
		if !ok || time.Until(deadline) > 5*time.Second {
			t.Error("handler does not have a bounded dependency context")
		}
	case <-time.After(time.Second):
		t.Fatal("request did not enter handler")
	}
	cancel()
	select {
	case err := <-response:
		if err != nil {
			t.Error(err)
		}
	case <-time.After(time.Second):
		t.Fatal("handler did not receive shutdown cancellation")
	}
}
