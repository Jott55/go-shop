package clog

import (
	"log"
	"os"
	"runtime"
)

type Level int

const (
	ERROR Level = iota
	INFO
	DEBUG
	HANDLER
	WARNING
)

func Log(level Level, info ...any) {
	Logger(level, 2, info...)
}

func Logger(level Level, skip int, info ...any) {

	var logger = log.New(os.Stdout, "Log: ", log.LstdFlags)

	switch level {
	case ERROR:
		logger.SetPrefix("ERROR: ")
		pc, filename, line, _ := runtime.Caller(skip)
		logger.Printf("at %s: %s %d\n\tWhat: %v", runtime.FuncForPC(pc).Name(), filename, line, info)
	case INFO:
		logger.SetPrefix("INFO: ")
		logger.Println(info...)
	case DEBUG:
		logger.SetFlags(0)
		logger.SetPrefix("DEBUG: ")
		logger.Println(info...)
	case WARNING:
		logger.SetPrefix("WARNING: ")
		logger.Println(info...)
	default:
		logger.Println(info...)
	}

}
