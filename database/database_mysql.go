package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	logger "github.com/Welasco/HubitatDeviceEvents/common/logger"
	"github.com/Welasco/HubitatDeviceEvents/model"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql_db struct {
	ConnectionString string
}

var db *sql.DB

func (mysqldb *Mysql_db) InitDB() (*sql.DB, error) {
	return sql.Open("mysql", mysqldb.ConnectionString)
}

func (mysqldb *Mysql_db) GetDevice(id int) (model.Device, error) {
	logger.Debug("[database][GetDevice] Getting device with id " + string(id))
	var device model.Device
	row := db.QueryRow("SELECT id, name, displayname, label FROM devices WHERE id = ?", id)
	err := row.Scan(&device.Id, &device.Name, &device.DisplayName, &device.Label)
	if err != nil {
		logger.Error("[database][GetDevice] Error getting device with id " + string(id))
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
	rows, err := db.Query("SELECT id, name, displayname, label FROM devices")
	if err != nil {
		logger.Error("[database][GetDevices] Error getting GetDevices")
		logger.Error("[database][GetDevices] Error: " + err.Error())
		return devices, err
	}

	for rows.Next() {
		var device model.Device
		err = rows.Scan(&device.Id, &device.Name, &device.DisplayName, &device.Label)
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
	_, err := db.Exec("INSERT INTO devices (name, displayname, label) VALUES (?, ?, ?)", device.Name, device.DisplayName, device.Label)
	if err != nil {
		logger.Error("[database][AddDevice] Error adding device")
		logger.Error("[database][AddDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][AddDevice] Adding device success")
	return err
}

func (mysqldb *Mysql_db) DeleteDevice(id int) error {
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
	_, err := db.Exec("UPDATE devices SET name = ?, displayname = ?, label = ? WHERE id = ?", device.Name, device.DisplayName, device.Label, device.Id)
	if err != nil {
		logger.Error("[database][UpdateDevice] Error updating device")
		logger.Error("[database][UpdateDevice] Error: " + err.Error())
		return err
	}
	logger.Debug("[database][UpdateDevice] Updating device success")
	return err
}

func (mysqldb *Mysql_db) CreateDB() error {

	// CREATE TABLE video_games
	// (
	// 	id    bigint unsigned not null primary key auto_increment,
	// 	name  VARCHAR(255)    NOT NULL,
	// 	displayname VARCHAR(255)    NOT NULL,
	// 	label  VARCHAR(255)         NOT NULL
	// );

	// root:password@tcp(127.0.0.1:3306)/hubitatdeviceevents
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
	log.Printf("rows affected: %d\n", no)
	logger.Info("[database][CreateDB] Rows affected during database creation: " + string(no))

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

	logger.Info("[database][CreateDB] Creating tables")
	res, err = dbCreate.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS devices (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, displayname VARCHAR(255) NOT NULL, label VARCHAR(255) NOT NULL)")
	if err != nil {
		logger.Error("[database][CreateDB] Error creating tables")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}
	logger.Info("[database][CreateDB] Creating tables success")

	no, err = res.RowsAffected()
	if err != nil {
		logger.Error("[database][CreateDB] Error fetching rows affected during table creation")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}

	logger.Info("[database][CreateDB] Closing database connection after table creation")
	dbCreate.Close()
	logger.Info("[database][CreateDB] Database full initialized success")
	return err
}

func (mysqldb *Mysql_db) GetDeviceEvents() error {
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
