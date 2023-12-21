package database

import (
	"context"
	"database/sql"
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
	var device model.Device
	row := db.QueryRow("SELECT id, name, displayname, label FROM devices WHERE id = ?", id)
	err := row.Scan(&device.Id, &device.Name, &device.DisplayName, &device.Label)
	if err != nil {
		return device, err
	}
	return device, nil
}

func (mysqldb *Mysql_db) GetDevices() ([]model.Device, error) {

	devices := []model.Device{}
	rows, err := db.Query("SELECT id, name, displayname, label FROM devices")
	if err != nil {
		return devices, err
	}

	for rows.Next() {
		var device model.Device
		err = rows.Scan(&device.Id, &device.Name, &device.DisplayName, &device.Label)
		if err != nil {
			return devices, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (mysqldb *Mysql_db) AddDevice(device *model.Device) error {
	_, err := db.Exec("INSERT INTO devices (name, displayname, label) VALUES (?, ?, ?)", device.Name, device.DisplayName, device.Label)
	return err
}

func (mysqldb *Mysql_db) DeleteDevice(id int) error {
	_, err := db.Exec("DELETE FROM devices WHERE id = ?", id)
	return err
}

func (mysqldb *Mysql_db) UpdateDevice(device *model.Device) error {
	_, err := db.Exec("UPDATE devices SET name = ?, displayname = ?, label = ? WHERE id = ?", device.Name, device.DisplayName, device.Label, device.Id)
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
	nodb_connectionstring := strings.SplitAfter(mysqldb.ConnectionString, "/")
	dbCreate, err := sql.Open("mysql", nodb_connectionstring[0])
	if err != nil {
		logger.Error("[database][CreateDB] Error initializing database")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := dbCreate.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+nodb_connectionstring[1])
	if err != nil {
		//log.Printf("Error %s when creating DB\n", err)
		//panic(err)
		logger.Error("[database][CreateDB] Error initializing database")
		logger.Error("[database][CreateDB] Error: " + err.Error())
		return err
	}

	err = dbCreate.PingContext(ctx)
	if err != nil {
		//log.Printf("Errors %s pinging DB", err)
		logger.Error("[database][CreateDB] Error initializing database PingContext")
		logger.Error("[database][CreateDB] Error: " + err.Error())
	}

	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return err
	}
	log.Printf("rows affected: %d\n", no)
	dbCreate.Close()

	dbCreate, err = sql.Open("mysql", mysqldb.ConnectionString)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return err
	}

	res, err = dbCreate.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS devices (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255) NOT NULL, displayname VARCHAR(255) NOT NULL, label VARCHAR(255) NOT NULL)")
	if err != nil {
		//log.Printf("Error %s when creating DB\n", err)
		return err
	}

	no, err = res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return err
	}
	log.Printf("rows affected: %d\n", no)
	dbCreate.Close()
	return err
}

func (mysqldb *Mysql_db) GetDeviceEvents() error {
	return nil
}

// func NewDB(ConnectionString string) (*Mysql_db, error) {
func NewDB(ConnectionString string) (*Mysql_db, error) {
	var mysqldb Mysql_db = Mysql_db{
		ConnectionString: ConnectionString,
	}
	var err error
	err = mysqldb.CreateDB()
	db, err = mysqldb.InitDB()
	return &mysqldb, err
}
