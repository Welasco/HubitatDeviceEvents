package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Welasco/HubitatDeviceEvents/model"
)

type testConnector struct {
	conn *testConnection
}

func (c testConnector) Connect(context.Context) (driver.Conn, error) { return c.conn, nil }
func (c testConnector) Driver() driver.Driver                        { return testDriver{} }

type testDriver struct{}

func (testDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("use testConnector")
}

type testConnection struct {
	exec  func(context.Context, string) (driver.Result, error)
	query func(context.Context, string) (driver.Rows, error)
}

func (*testConnection) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("unexpected Prepare")
}
func (*testConnection) Begin() (driver.Tx, error) {
	return nil, errors.New("unexpected Begin")
}
func (*testConnection) Close() error { return nil }
func (c *testConnection) ExecContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Result, error) {
	return c.exec(ctx, query)
}
func (c *testConnection) QueryContext(ctx context.Context, query string, _ []driver.NamedValue) (driver.Rows, error) {
	return c.query(ctx, query)
}

type testRows struct {
	columns []string
	values  [][]driver.Value
	err     error
	closed  atomic.Bool
}

func (r *testRows) Columns() []string { return r.columns }
func (r *testRows) Close() error {
	r.closed.Store(true)
	return nil
}
func (r *testRows) Next(dest []driver.Value) error {
	if len(r.values) == 0 {
		if r.err != nil {
			return r.err
		}
		return io.EOF
	}
	copy(dest, r.values[0])
	r.values = r.values[1:]
	return nil
}

func testDatabase(t *testing.T, conn *testConnection) *Mysql_db {
	t.Helper()
	db := sql.OpenDB(testConnector{conn})
	db.SetMaxOpenConns(1)
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	})
	return &Mysql_db{db: db}
}

func TestDatabaseOperationsHonorRequestDeadline(t *testing.T) {
	methods := map[string]func(*Mysql_db, context.Context) error{
		"GetDevice": func(db *Mysql_db, ctx context.Context) error {
			_, err := db.GetDevice(ctx, "42")
			return err
		},
		"GetDevices": func(db *Mysql_db, ctx context.Context) error {
			_, err := db.GetDevices(ctx)
			return err
		},
		"AddDevice": func(db *Mysql_db, ctx context.Context) error {
			return db.AddDevice(ctx, &model.Device{Id: "42"})
		},
		"DeleteDevice": func(db *Mysql_db, ctx context.Context) error {
			return db.DeleteDevice(ctx, "42")
		},
		"UpdateDevice": func(db *Mysql_db, ctx context.Context) error {
			return db.UpdateDevice(ctx, &model.Device{Id: "42"})
		},
		"GetDeviceEvents": func(db *Mysql_db, ctx context.Context) error {
			_, err := db.GetDeviceEvents(ctx)
			return err
		},
		"GetDeviceEventId": func(db *Mysql_db, ctx context.Context) error {
			_, err := db.GetDeviceEventId(ctx, "42", model.UrlQueries{})
			return err
		},
		"RegisterDeviceEvent": func(db *Mysql_db, ctx context.Context) error {
			return db.RegisterDeviceEvent(ctx, &model.DeviceEvent{DeviceId: "42"})
		},
	}
	for name, method := range methods {
		t.Run(name, func(t *testing.T) {
			wait := func(ctx context.Context) error {
				<-ctx.Done()
				return ctx.Err()
			}
			db := testDatabase(t, &testConnection{
				exec: func(ctx context.Context, _ string) (driver.Result, error) {
					return nil, wait(ctx)
				},
				query: func(ctx context.Context, _ string) (driver.Rows, error) {
					return nil, wait(ctx)
				},
			})
			ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
			defer cancel()
			start := time.Now()
			if err := method(db, ctx); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("expected deadline error, got %v", err)
			}
			if elapsed := time.Since(start); elapsed > time.Second {
				t.Errorf("operation ignored request deadline: %s", elapsed)
			}
		})
	}
}

func TestQueryErrorsReleaseRowsAndAreNotPartialSuccess(t *testing.T) {
	rowErr := errors.New("connection lost while reading rows")
	for _, name := range []string{"GetDevices", "GetDeviceEvents", "GetDeviceEventId", "scan error"} {
		t.Run(name, func(t *testing.T) {
			rows := &testRows{columns: []string{"id", "name", "label", "type", "room"}, err: rowErr}
			if name == "scan error" {
				rows.values = [][]driver.Value{{nil, "name", "label", "type", "room"}}
				rows.err = nil
			}
			db := testDatabase(t, &testConnection{
				query: func(context.Context, string) (driver.Rows, error) { return rows, nil },
			})
			var err error
			switch name {
			case "GetDeviceEvents":
				_, err = db.GetDeviceEvents(context.Background())
			case "GetDeviceEventId":
				_, err = db.GetDeviceEventId(context.Background(), "42", model.UrlQueries{})
			default:
				_, err = db.GetDevices(context.Background())
			}
			if err == nil || (name != "scan error" && !errors.Is(err, rowErr)) {
				t.Errorf("expected row failure, got %v", err)
			}
			if !rows.closed.Load() || db.db.Stats().InUse != 0 {
				t.Error("failed query retained an open result set or connection")
			}
		})
	}
}

func TestRegisterDeviceEventPreservesAutoRegistrationAndInsertFailure(t *testing.T) {
	insertErr := errors.New("event insert failed")
	for _, exists := range []bool{false, true} {
		name := "unknown device"
		if exists {
			name = "existing device"
		}
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/42" {
					t.Errorf("unexpected Hubitat path: %s", r.URL.Path)
				}
				io.WriteString(w, `{"Id":"42","name":"sensor","label":"sensor","type":"sensor","room":"room"}`)
			}))
			defer server.Close()
			t.Setenv("HubitatGetDevicesUrl", server.URL+"?access_token=test")
			eventAttempts := 0
			deviceAdded := false
			db := testDatabase(t, &testConnection{
				exec: func(ctx context.Context, query string) (driver.Result, error) {
					if strings.HasPrefix(query, "INSERT INTO devices ") {
						deviceAdded = true
						return driver.RowsAffected(1), nil
					}
					eventAttempts++
					if eventAttempts == 1 {
						return nil, insertErr
					}
					if !deviceAdded {
						t.Error("retried event before registering its device")
					}
					return driver.RowsAffected(1), nil
				},
				query: func(context.Context, string) (driver.Rows, error) {
					rows := &testRows{columns: []string{"id", "name", "label", "type", "room"}}
					if exists {
						rows.values = [][]driver.Value{{"42", "name", "label", "type", "room"}}
					}
					return rows, nil
				},
			})
			err := db.RegisterDeviceEvent(context.Background(), &model.DeviceEvent{DeviceId: "42"})
			if exists {
				if !errors.Is(err, insertErr) || eventAttempts != 1 || deviceAdded {
					t.Fatalf("failed write was acknowledged or retried unexpectedly: %v", err)
				}
			} else if err != nil || eventAttempts != 2 || !deviceAdded {
				t.Fatalf("device/event auto-registration failed: %v", err)
			}
		})
	}
}

func TestDatabaseStartupCancelsStalledMySQLHandshake(t *testing.T) {
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			accepted <- conn
		}
		close(accepted)
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	db, err := NewMySQLDB(ctx, "test@tcp("+listener.Addr().String()+")/events")
	listener.Close()
	if conn := <-accepted; conn != nil {
		conn.Close()
	} else {
		t.Fatal("test did not establish a MySQL connection")
	}
	if db != nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected canceled database startup, got %v, %v", db, err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("database startup ignored cancellation: %s", elapsed)
	}
}
