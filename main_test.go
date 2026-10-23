package main

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/transport"
	"github.com/gofiber/fiber/v2"
)

func TestShutdownDrainsEventBeforeCancelingAndClosingDatabase(t *testing.T) {
	requestContext, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()
	app := transport.Setup(requestContext)
	entered := make(chan struct{})
	release := make(chan struct{})
	var saved atomic.Bool
	app.Post("/test-event", func(c *fiber.Ctx) error {
		close(entered)
		select {
		case <-release:
			saved.Store(true)
			return c.SendString("saved")
		case <-c.UserContext().Done():
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}
	})
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- app.Listener(listener) }()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := app.ShutdownWithContext(ctx); err != nil {
			t.Error(err)
		}
	}()

	response := make(chan error, 1)
	go func() {
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Post("http://"+listener.Addr().String()+"/test-event", "application/json", strings.NewReader("{}"))
		if err == nil {
			defer resp.Body.Close()
			var body []byte
			body, err = io.ReadAll(resp.Body)
			if err == nil && (resp.StatusCode != http.StatusOK || string(body) != "saved") {
				err = errors.New("event was not acknowledged after saving")
			}
		}
		response <- err
	}()
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("request did not enter handler")
	}

	draining := make(chan struct{})
	finished := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	go func() {
		finished <- shutdown(ctx, func(ctx context.Context) error {
			close(draining)
			return app.ShutdownWithContext(ctx)
		}, cancelRequests, func() error {
			if !saved.Load() {
				return errors.New("database closed before event was saved")
			}
			if requestContext.Err() == nil {
				return errors.New("requests were not canceled before database close")
			}
			return nil
		})
	}()
	<-draining
	if requestContext.Err() != nil {
		t.Error("in-flight request was canceled before drain finished")
	}
	close(release)
	if err := <-response; err != nil {
		t.Error(err)
	}
	if err := <-finished; err != nil {
		t.Error(err)
	}
	if err := <-serverDone; err != nil {
		t.Error(err)
	}
}

func TestShutdownDeadlineBoundsBlockedDrainAndCleanup(t *testing.T) {
	for _, stage := range []string{"drain", "cleanup"} {
		t.Run(stage, func(t *testing.T) {
			requestContext, cancelRequests := context.WithCancel(context.Background())
			defer cancelRequests()
			release := make(chan struct{})
			workerDone := make(chan struct{})
			defer func() {
				close(release)
				select {
				case <-workerDone:
				case <-time.After(time.Second):
					t.Error("shutdown worker did not finish")
				}
			}()
			drain := func(context.Context) error {
				if stage == "drain" {
					<-release
				}
				return nil
			}
			closeDB := func() error {
				defer close(workerDone)
				if stage == "cleanup" {
					<-release
				}
				return nil
			}
			ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
			defer cancel()
			start := time.Now()
			err := shutdown(ctx, drain, cancelRequests, closeDB)
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected deadline error, got %v", err)
			}
			if elapsed := time.Since(start); elapsed > time.Second {
				t.Errorf("shutdown exceeded its exit budget: %s", elapsed)
			}
			if !errors.Is(requestContext.Err(), context.Canceled) {
				t.Error("remaining requests were not canceled")
			}
		})
	}
}

func TestShutdownCancelsAfterDrainFailureAndReportsCleanupFailure(t *testing.T) {
	requestContext, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()
	drainErr := errors.New("drain failed")
	closeErr := errors.New("database close failed")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := shutdown(ctx, func(context.Context) error {
		return drainErr
	}, cancelRequests, func() error {
		if requestContext.Err() == nil {
			t.Error("database cleanup started without canceling requests")
		}
		return closeErr
	})
	if !errors.Is(err, drainErr) || !errors.Is(err, closeErr) {
		t.Fatalf("shutdown lost errors: %v", err)
	}
}

func TestShutdownDeadlineIncludesBlockedLogging(t *testing.T) {
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	original := os.Stdout
	os.Stdout = writer
	err = logger.Init("", "Inf")
	os.Stdout = original
	if err != nil {
		t.Fatal(err)
	}
	requestContext, cancelRequests := context.WithCancel(context.Background())
	defer cancelRequests()
	workerDone := make(chan struct{})
	defer func() {
		reader.Close()
		select {
		case <-workerDone:
		case <-time.After(time.Second):
			t.Error("logging did not unblock after closing pipe")
		}
		writer.Close()
		if err := logger.Init("", "Err"); err != nil {
			t.Error(err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	start := time.Now()
	err = shutdown(ctx, func(context.Context) error {
		logger.Info("[test][shutdown] " + strings.Repeat("x", 1024*1024))
		return nil
	}, cancelRequests, func() error {
		close(workerDone)
		return nil
	})
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(start) > time.Second {
		t.Fatalf("blocked logging escaped shutdown deadline: %v", err)
	}
	if requestContext.Err() == nil {
		t.Error("blocked logging prevented request cancellation")
	}
}
