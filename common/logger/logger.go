package logger

import (
	"io"
	"log"
	"os"
)

var l = log.New(os.Stdout, "", log.LstdFlags)

var level int

// array of log levels to match the config file and the const
var loglevelstr = []string{"Err", "Warn", "Inf", "Deb"}

const (
	// ErrorLevel level. Used for errors that should definitely be noted.
	Err = iota
	// WarningLevel level. Non-critical entries that deserve eyes.
	Warn
	// InfoLevel level. General operational entries about what's going on inside the application.
	Inf
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	Deb
)

// Init initializes the logger
func Init(file string, llevel string) error {
	level = Err
	for i, v := range loglevelstr {
		if v == llevel {
			level = i
		}
	}

	l.SetOutput(os.Stdout)
	if file != "" {
		logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		l.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
	l.SetFlags(log.LstdFlags)
	return nil
}

func Error(msg ...interface{}) {
	if level >= Err {
		l.Println("ERROR:", msg)
	}
}

func Warning(msg ...interface{}) {
	if level >= Warn {
		l.Println("WARNING:", msg)
	}
}

func Info(msg ...interface{}) {
	if level >= Inf {
		l.Println("INFO:", msg)
	}
}

func Debug(msg ...interface{}) {
	if level >= Deb {
		l.Println("DEBUG:", msg)
	}
}
