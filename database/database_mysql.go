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
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Mysql_db struct {
	ConnectionString string
}

var db *sql.DB

func (mysqldb *Mysql_db) InitDB() (*sql.DB, error) {
	return sql.Open("mysql", mysqldb.ConnectionString)
}

func (mysqldb *Mysql_db) GetDevice(id string) (model.Device, error) {
	logger.Debug("[database][GetDevice] Getting device with id " + id)
	var device model.Device
	row := db.QueryRow("SELECT id, name, label, type, room FROM devices WHERE id = ?", id)
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

func (mysqldb *Mysql_db) GetDevices() ([]model.Device, error) {
	logger.Debug("[database][GetDevices] Getting devices")
	devices := []model.Device{}
	rows, err := db.Query("SELECT id, name, label, type, room FROM devices")
	if err != nil {
		logger.Error("[database][GetDevices] Error getting GetDevices")
		logger.Error("[database][GetDevices] Error: " + err.Error())
		return devices, err
	}

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
	json, _ := json.Marshal(devices)
	logger.Debug("[database][GetDevices] Getting devices: " + string(json))
	return devices, nil
}

func (mysqldb *Mysql_db) AddDevice(device *model.Device) error {
	logger.Debug("[database][AddDevice] Adding device")
	_, err := db.Exec("INSERT INTO devices (id, name, label, type, room) VALUES (?, ?, ?, ?, ?)", device.Id, device.Name, device.Label, device.Type, device.Room)
	if err != nil {
		logger.Error("[database][AddDevice] Error adding device")
		logger.Error("[database][AddDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][AddDevice] Adding device success")
	return err
}

func (mysqldb *Mysql_db) DeleteDevice(id string) error {
	logger.Debug("[database][DeleteDevice] Deleting device")
	_, err := db.Exec("DELETE FROM devices WHERE id = ?", id)
	if err != nil {
		logger.Error("[database][DeleteDevice] Error deleting device")
		logger.Error("[database][DeleteDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][DeleteDevice] Deleting device success")
	return err
}

func (mysqldb *Mysql_db) UpdateDevice(device *model.Device) error {
	logger.Debug("[database][UpdateDevice] Updating device")
	_, err := db.Exec("UPDATE devices SET name = ?, label = ?, type = ?, room = ? WHERE id = ?", device.Name, device.Label, device.Type, device.Room, device.Id)
	if err != nil {
		logger.Error("[database][UpdateDevice] Error updating device")
		logger.Error("[database][UpdateDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][UpdateDevice] Updating device success")
	return err
}

func (mysqldb *Mysql_db) CreateDB() error {

	logger.Info("[database][CreateDB] Initializing database connection to create database")
	nodb_connectionstring := strings.SplitAfter(mysqldb.ConnectionString, "/")
	dbCreate, err := sql.Open("mysql", nodb_connectionstring[0])
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database for database creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}

	logger.Info("[database][CreateDB] Creating Database " + nodb_connectionstring[1])
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := dbCreate.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+nodb_connectionstring[1])
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Database creation success " + nodb_connectionstring[1])

	logger.Info("[database][CreateDB] Ping Database")
	err = dbCreate.PingContext(ctx)
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database PingContext")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Ping Database success")

	logger.Info("[database][CreateDB] Checking affected rows during database creation")

	no, err := res.RowsAffected()
	if err != nil {
		logger.Error("[database][CreateDB] Error fetching rows affected during database creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Rows affected during database creation: " + fmt.Sprint(no))

	logger.Info("[database][CreateDB] Closing database connection after database creation")
	dbCreate.Close()

	logger.Info("[database][CreateDB] Initializing database connection to create tables")
	dbCreate, err = sql.Open("mysql", mysqldb.ConnectionString)
	if err != nil {
		logger.Error("[database][CreateDB] Error connecting to database for table creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
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
	res, err = dbCreate.ExecContext(ctx, create_table_devices)
	if err != nil {
		logger.Error("[database][CreateDB] Error creating tables")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Creating table devices success")

	no, err = res.RowsAffected()
	if err != nil {
		logger.Error("[database][CreateDB] Error fetching rows affected during table devices creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Rows affected during table devices creation: " + fmt.Sprint(no))

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
	res, err = dbCreate.ExecContext(ctx, create_table_deviceevents)
	if err != nil {
		logger.Error("[database][CreateDB] Error creating table deviceevents")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Creating table deviceevents success")

	no, err = res.RowsAffected()
	if err != nil {
		logger.Error("[database][CreateDB] Error fetching rows affected during table deviceevents creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Rows affected during table deviceevents creation: " + fmt.Sprint(no))

	logger.Info("[database][CreateDB] Closing database connection after table creation")
	dbCreate.Close()
	logger.Info("[database][CreateDB] Database full initialized success")
	return err
}

func (mysqldb *Mysql_db) GetDeviceEvents() ([]model.DeviceEvent, error) {
	logger.Debug("[database][GetDeviceEvents] Getting devices events")
	deviceevents := []model.DeviceEvent{}
	rows, err := db.Query("SELECT timestamp, name, value, displayName, deviceId, descriptionText, unit, type, data FROM deviceevents")
	if err != nil {
		logger.Error("[database][GetDeviceEvents] Error getting devices events")
		logger.Error("[database][GetDeviceEvents] Error: " + err.Error())
		return deviceevents, err
	}

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
	json, _ := json.Marshal(deviceevents)
	logger.Debug("[database][GetDeviceEvents] Getting device events: " + string(json))
	return deviceevents, nil
}

func (mysqldb *Mysql_db) GetDeviceEventId(id string) ([]model.DeviceEvent, error) {
	logger.Debug("[database][GetDeviceEventId] Getting events from device with id " + id)
	deviceevents := []model.DeviceEvent{}
	rows, err := db.Query("SELECT timestamp, name, value, displayName, deviceId, descriptionText, unit, type, data FROM deviceevents WHERE deviceId = ?", id)
	if err != nil {
		logger.Error("[database][GetDeviceEventId] Error getting events from device with id " + id)
		logger.Error("[database][GetDeviceEventId] Error: " + err.Error())
		return deviceevents, err
	}

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
	json, _ := json.Marshal(deviceevents)
	logger.Debug("[database][GetDeviceEventId] Getting events from device with id " + id + ": " + string(json))
	return deviceevents, nil
}

func (mysqldb *Mysql_db) RegisterDeviceEvent(event *model.DeviceEvent) error {
	logger.Debug("[database][RegisterDeviceEvent] Registering device event")
	_, err := db.Exec("INSERT INTO deviceevents (name, value, displayName, deviceId, descriptionText, unit, type, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", event.Name, event.Value, event.DisplayName, event.DeviceId, event.DescriptionText, event.Unit, event.Type, event.Data)
	if err != nil {
		logger.Error("[database][RegisterDeviceEvent] Error registering device")
		logger.Error("[database][RegisterDeviceEvent] Error: " + err.Error())
		//return err

		// Check if device exists
		device, err := mysqldb.GetDevice(event.DeviceId)
		if err != nil {
			logger.Error("[device][RegisterDeviceEvent] Error reading devices from database")
			logger.Error("[device][RegisterDeviceEvent] Error: Device " + event.DeviceId + " not found")
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
			hubitatUrl_builder := strings.Split(hubitatUrl, "?")
			respbody, err := http.Http_client(hubitatUrl_builder[0] + "/" + event.DeviceId + "?" + hubitatUrl_builder[1])
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
			err = mysqldb.AddDevice(&device)
			if err != nil {
				logger.Error("[device][RegisterDeviceEvent] Error adding device")
				logger.Error("[device][RegisterDeviceEvent] Error: " + err.Error())
				return err
			}
			logger.Debug("[database][RegisterDeviceEvent] Device added")

			logger.Debug("[database][RegisterDeviceEvent] Registering device event after registering device")
			_, err = db.Exec("INSERT INTO deviceevents (name, value, displayName, deviceId, descriptionText, unit, type, data) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", event.Name, event.Value, event.DisplayName, event.DeviceId, event.DescriptionText, event.Unit, event.Type, event.Data)
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

func NewMySQLDB(ConnectionString string) (*Mysql_db, error) {
	var mysqldb Mysql_db = Mysql_db{
		ConnectionString: ConnectionString,
	}
	var err error
	logger.Info("[database][NewMySQLDB] Invoking database creation")
	err = mysqldb.CreateDB()
	logger.Info("[database][NewMySQLDB] Invoking database initialization")
	db, err = mysqldb.InitDB()
	if err != nil {
		logger.Error("[database][NewMySQLDB] Error initializing database connection")
		logger.Error("[database][NewMySQLDB] Error: " + err.Error())
	}
	return &mysqldb, err
}
