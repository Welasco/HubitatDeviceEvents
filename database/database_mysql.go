package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Welasco/HubitatDeviceEvents/common/http"
	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/model"
	"github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Mysql_db struct {
	ConnectionString string
	db               *sql.DB
}

const operationTimeout = 5 * time.Second

func (mysqldb *Mysql_db) InitDB() (*sql.DB, error) {
	return sql.Open("mysql", mysqldb.ConnectionString)
}

func (mysqldb *Mysql_db) Close() error {
	if mysqldb.db == nil {
		return nil
	}
	logger.Info("[database][Close] Closing database connection")
	return mysqldb.db.Close()
}

func (mysqldb *Mysql_db) GetDevice(ctx context.Context, id string) (model.Device, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][GetDevice] Getting device with id " + id)
	var device model.Device
	row := mysqldb.db.QueryRowContext(ctx, "SELECT id, name, label, type, room FROM devices WHERE id = ?", id)
	err := row.Scan(&device.Id, &device.Name, &device.Label, &device.Type, &device.Room)
	if err != nil {
		logger.Error("[database][GetDevice] Error getting device with id " + id)
		logger.Error("[database][GetDevice] Error: " + err.Error())
		return device, err
	}
	json, _ := json.Marshal(device)
	logger.Debug("[database][GetDevice] Getting device with id " + string(json))
	return device, nil
}

func (mysqldb *Mysql_db) GetDevices(ctx context.Context) ([]model.Device, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][GetDevices] Getting devices")
	devices := []model.Device{}
	rows, err := mysqldb.db.QueryContext(ctx, "SELECT id, name, label, type, room FROM devices order by label")
	if err != nil {
		logger.Error("[database][GetDevices] Error getting GetDevices")
		logger.Error("[database][GetDevices] Error: " + err.Error())
		return devices, err
	}
	defer rows.Close()

	for rows.Next() {
		var device model.Device
		err = rows.Scan(&device.Id, &device.Name, &device.Label, &device.Type, &device.Room)
		if err != nil {
			logger.Error("[database][GetDevices] Error scanning GetDevices rows")
			logger.Error("[database][GetDevices] Error: " + err.Error())
			return devices, err
		}
		devices = append(devices, device)
	}
	if err := rows.Err(); err != nil {
		logger.Error("[database][GetDevices] Error reading rows: " + err.Error())
		return devices, err
	}
	json, _ := json.Marshal(devices)
	logger.Debug("[database][GetDevices] Getting devices: " + string(json))
	return devices, nil
}

func (mysqldb *Mysql_db) AddDevice(ctx context.Context, device *model.Device) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][AddDevice] Adding device")
	_, err := mysqldb.db.ExecContext(ctx, "INSERT INTO devices (id, name, label, type, room) VALUES (?, ?, ?, ?, ?)", device.Id, device.Name, device.Label, device.Type, device.Room)
	if err != nil {
		logger.Error("[database][AddDevice] Error adding device")
		logger.Error("[database][AddDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][AddDevice] Adding device success")
	return err
}

func (mysqldb *Mysql_db) DeleteDevice(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][DeleteDevice] Deleting device")
	_, err := mysqldb.db.ExecContext(ctx, "DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		logger.Error("[database][DeleteDevice] Error deleting device")
		logger.Error("[database][DeleteDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][DeleteDevice] Deleting device success")
	return err
}

func (mysqldb *Mysql_db) UpdateDevice(ctx context.Context, device *model.Device) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][UpdateDevice] Updating device")
	_, err := mysqldb.db.ExecContext(ctx, "UPDATE devices SET name = ?, label = ?, type = ?, room = ? WHERE id = ?", device.Name, device.Label, device.Type, device.Room, device.Id)
	if err != nil {
		logger.Error("[database][UpdateDevice] Error updating device")
		logger.Error("[database][UpdateDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][UpdateDevice] Updating device success")
	return err
}

func (mysqldb *Mysql_db) CreateDB(ctx context.Context) (err error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	cfg, err := mysql.ParseDSN(mysqldb.ConnectionString)
	if err != nil {
		return err
	}
	databaseName := cfg.DBName
	if databaseName == "" {
		return errors.New("database name is required in ConnectionString")
	}
	cfg.DBName = ""
	logger.Info("[database][CreateDB] Initializing database connection to create database")
	dbCreate, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database for database creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	defer func() {
		err = errors.Join(err, dbCreate.Close())
	}()
	_, err = dbCreate.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS `"+strings.ReplaceAll(databaseName, "`", "``")+"`")
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Initializing database connection to create tables")
	dbTables, err := sql.Open("mysql", mysqldb.ConnectionString)
	if err != nil {
		logger.Error("[database][CreateDB] Error connecting to database for table creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	defer func() {
		err = errors.Join(err, dbTables.Close())
	}()
	logger.Info("[database][CreateDB] Connected to database for table creation")

	logger.Info("[database][CreateDB] Creating table devices")
	var create_table_devices string = `
	CREATE TABLE IF NOT EXISTS devices (
		id BIGINT NOT NULL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		label VARCHAR(255) NOT NULL,
		type VARCHAR(255) NOT NULL,
		room VARCHAR(255) NULL
	);
	`
	_, err = dbTables.ExecContext(ctx, create_table_devices)
	if err != nil {
		logger.Error("[database][CreateDB] Error creating tables")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Creating table devices success")

	logger.Info("[database][CreateDB] Creating deviceevents devices")
	var create_table_deviceevents string = `
	CREATE TABLE IF NOT EXISTS deviceevents (
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		name VARCHAR(255) NOT NULL,
		value VARCHAR(255) NOT NULL,
		displayName VARCHAR(255) NOT NULL,
		deviceId BIGINT NOT NULL,
		descriptionText VARCHAR(255),
		unit VARCHAR(255) NULL,
		type VARCHAR(255) NULL,
		data VARCHAR(255) NULL,
		CONSTRAINT deviceId FOREIGN KEY (deviceId) REFERENCES devices (id)
	);
	`
	_, err = dbTables.ExecContext(ctx, create_table_deviceevents)
	if err != nil {
		logger.Error("[database][CreateDB] Error creating table deviceevents")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Creating table deviceevents success")

	logger.Info("[database][CreateDB] Database full initialized success")
	return err
}

func (mysqldb *Mysql_db) GetDeviceEvents(ctx context.Context) ([]model.DeviceEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][GetDeviceEvents] Getting devices events")
	deviceevents := []model.DeviceEvent{}
	rows, err := mysqldb.db.QueryContext(ctx, "SELECT timestamp, name, value, displayName, deviceId, descriptionText, unit, type, data FROM deviceevents")
	if err != nil {
		logger.Error("[database][GetDeviceEvents] Error getting devices events")
		logger.Error("[database][GetDeviceEvents] Error: " + err.Error())
		return deviceevents, err
	}
	defer rows.Close()

	for rows.Next() {
		var deviceevent model.DeviceEvent
		err = rows.Scan(&deviceevent.TimeStamp, &deviceevent.Name, &deviceevent.Value, &deviceevent.DisplayName, &deviceevent.DeviceId, &deviceevent.DescriptionText, &deviceevent.Unit, &deviceevent.Type, &deviceevent.Data)
		if err != nil {
			logger.Error("[database][GetDeviceEvents] Error scanning device events rows")
			logger.Error("[database][GetDeviceEvents] Error: " + err.Error())
			return deviceevents, err
		}
		deviceevents = append(deviceevents, deviceevent)
	}
	if err := rows.Err(); err != nil {
		logger.Error("[database][GetDeviceEvents] Error reading rows: " + err.Error())
		return deviceevents, err
	}
	json, _ := json.Marshal(deviceevents)
	logger.Debug("[database][GetDeviceEvents] Getting device events: " + string(json))
	return deviceevents, nil
}

func (mysqldb *Mysql_db) GetDeviceEventId(ctx context.Context, id string, timeStamp model.UrlQueries) ([]model.DeviceEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][GetDeviceEventId] Getting events from device with id " + id)
	deviceevents := []model.DeviceEvent{}
	var rows *sql.Rows
	var err error

	if len(timeStamp.Start) > 0 && len(timeStamp.End) > 0 {
		rows, err = mysqldb.db.QueryContext(ctx, "SELECT timestamp, name, value, displayName, deviceId, descriptionText, unit, type, data FROM deviceevents WHERE deviceId = ? and timestamp BETWEEN ? and ?", id, timeStamp.Start, timeStamp.End)
		if err != nil {
			logger.Error("[database][GetDeviceEventId] Error getting events from device with id " + id)
			logger.Error("[database][GetDeviceEventId] Error: " + err.Error())
			return deviceevents, err
		}
	} else {
		rows, err = mysqldb.db.QueryContext(ctx, "SELECT timestamp, name, value, displayName, deviceId, descriptionText, unit, type, data FROM deviceevents WHERE deviceId = ?", id)
		if err != nil {
			logger.Error("[database][GetDeviceEventId] Error getting events from device with id " + id)
			logger.Error("[database][GetDeviceEventId] Error: " + err.Error())
			return deviceevents, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var deviceevent model.DeviceEvent
		err = rows.Scan(&deviceevent.TimeStamp, &deviceevent.Name, &deviceevent.Value, &deviceevent.DisplayName, &deviceevent.DeviceId, &deviceevent.DescriptionText, &deviceevent.Unit, &deviceevent.Type, &deviceevent.Data)
		if err != nil {
			logger.Error("[database][GetDeviceEventId] Error scanning events from device with id " + id + " rows")
			logger.Error("[database][GetDeviceEventId] Error: " + err.Error())
			return deviceevents, err
		}
		deviceevents = append(deviceevents, deviceevent)
	}
	if err := rows.Err(); err != nil {
		logger.Error("[database][GetDeviceEventId] Error reading rows: " + err.Error())
		return deviceevents, err
	}
	json, _ := json.Marshal(deviceevents)
	logger.Debug("[database][GetDeviceEventId] Getting events from device with id " + id + ": " + string(json))
	return deviceevents, nil
}

func (mysqldb *Mysql_db) RegisterDeviceEvent(ctx context.Context, event *model.DeviceEvent) error {
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	logger.Debug("[database][RegisterDeviceEvent] Registering device event")
	_, err := mysqldb.db.ExecContext(ctx, "INSERT INTO deviceevents (name, value, displayName, deviceId, descriptionText, unit, type, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", event.Name, event.Value, event.DisplayName, event.DeviceId, event.DescriptionText, event.Unit, event.Type, event.Data)
	if err != nil {
		logger.Error("[database][RegisterDeviceEvent] Error registering device")
		logger.Error("[database][RegisterDeviceEvent] Error: " + err.Error())
		insertErr := err
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Check if device exists
		device, err := mysqldb.GetDevice(ctx, event.DeviceId)
		if err != nil {
			logger.Error("[device][RegisterDeviceEvent] Error reading devices from database")
			logger.Error("[device][RegisterDeviceEvent] Error: Device " + event.DeviceId + " not found")
			if !errors.Is(err, sql.ErrNoRows) {
				return err
			}
		}
		if device.Id != "" {
			return insertErr
		}
		if device.Id == "" {
			logger.Debug("[database][RegisterDeviceEvent] Device not found, registering device")

			logger.Debug("[database][RegisterDeviceEvent] Retrieving device from Hubitat Hub")
			hubitatUrl := os.Getenv("HubitatGetDevicesUrl")
			if hubitatUrl == "" {
				logger.Error("[database][RegisterDeviceEvent] Error getting HubitatGetDevicesUrl from environment variable or .env file")
				logger.Error("[database][RegisterDeviceEvent] Unable to register device")
				return errors.New("error getting HubitatGetDevicesUrl from environment variable or .env file")
			}
			hubitatUrl_builder := strings.SplitN(hubitatUrl, "?", 2)
			if len(hubitatUrl_builder) != 2 {
				return errors.New("HubitatGetDevicesUrl must include an access token query")
			}
			respbody, err := http.Http_client(ctx, hubitatUrl_builder[0]+"/"+event.DeviceId+"?"+hubitatUrl_builder[1])
			if err != nil {
				logger.Error("[database][RegisterDeviceEvent] Failed to connect to Hubitat Hub")
				logger.Error("[database][RegisterDeviceEvent] Error: " + err.Error())
				logger.Error("[database][RegisterDeviceEvent] Unable to register device")
				return err
			}
			logger.Debug("[database][RegisterDeviceEvent] Device retrieved from Hubitat Hub")
			var device model.Device
			err = json.Unmarshal([]byte(respbody), &device)
			if err != nil {
				logger.Error("[database][RegisterDeviceEvent] Error unmarshalling device")
				logger.Error("[database][RegisterDeviceEvent] Error: " + err.Error())
				logger.Error("[database][RegisterDeviceEvent] Unable to register device")
				return err
			}
			logger.Debug("[database][RegisterDeviceEvent] Device unmarshalled")

			logger.Debug("[database][RegisterDeviceEvent] Adding device")
			err = mysqldb.AddDevice(ctx, &device)
			if err != nil {
				logger.Error("[device][RegisterDeviceEvent] Error adding device")
				logger.Error("[device][RegisterDeviceEvent] Error: " + err.Error())
				return err
			}
			logger.Debug("[database][RegisterDeviceEvent] Device added")

			logger.Debug("[database][RegisterDeviceEvent] Registering device event after registering device")
			_, err = mysqldb.db.ExecContext(ctx, "INSERT INTO deviceevents (name, value, displayName, deviceId, descriptionText, unit, type, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", event.Name, event.Value, event.DisplayName, event.DeviceId, event.DescriptionText, event.Unit, event.Type, event.Data)
			if err != nil {
				logger.Error("[database][RegisterDeviceEvent] Error registering device to database after registering device")
				logger.Error("[database][RegisterDeviceEvent] Error: " + err.Error())
				return err
			}
			logger.Debug("[database][RegisterDeviceEvent] Registering device event after registering device success")
		}
	}

	logger.Debug("[database][RegisterDeviceEvent] Registering event device success")
	return nil
}

func NewMySQLDB(ctx context.Context, ConnectionString string) (*Mysql_db, error) {
	cfg, err := mysql.ParseDSN(ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("invalid MySQL ConnectionString: %w", err)
	}
	if cfg.Timeout <= 0 || cfg.Timeout > operationTimeout {
		cfg.Timeout = operationTimeout
	}
	if cfg.ReadTimeout <= 0 || cfg.ReadTimeout > operationTimeout {
		cfg.ReadTimeout = operationTimeout
	}
	if cfg.WriteTimeout <= 0 || cfg.WriteTimeout > operationTimeout {
		cfg.WriteTimeout = operationTimeout
	}
	var mysqldb Mysql_db = Mysql_db{
		ConnectionString: cfg.FormatDSN(),
	}
	logger.Info("[database][NewMySQLDB] Invoking database creation")
	if err := mysqldb.CreateDB(ctx); err != nil {
		return nil, err
	}
	logger.Info("[database][NewMySQLDB] Invoking database initialization")
	mysqldb.db, err = mysqldb.InitDB()
	if err != nil {
		logger.Error("[database][NewMySQLDB] Error initializing database connection")
		logger.Error("[database][NewMySQLDB] Error: " + err.Error())
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	if err := mysqldb.db.PingContext(ctx); err != nil {
		return nil, errors.Join(err, mysqldb.Close())
	}
	return &mysqldb, nil
}
