// Package logger provides a leveled logging solution.
// It supports three levels of logging: INFO, ERROR, and FATAL.
// Each level has its own logger, and each logger writes to its own output.
// All logs are written to LOG_FILE and console,
// and error logs are also written to ERROR_LOG_FILE.
//
// Usage:
// Before using the logger, create a new logger by calling the `NewLogger` function,
// passing in the paths to the log file and error log file, and the desired log level.
// Then, use the `Info`, `Error`, and `Fatal` methods to log messages at the respective levels.
// Finally, call the `Close` method to close the log files when you're done with the logger.
//
// Example:
// Create a new logger:
// logger, err := logger.NewLogger("log.txt", "error_log.txt", logger.INFO)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Log an info message:
// logger.Info("This is an info message")
//
// Log an error message: (This will also be written to the error log file)
// logger.Error("This is an error message")
//
// Log a fatal message: (Same as error, but keep in mind that this will exit the program)
// logger.Fatal("This is a fatal message")
//
// Close the logger:
// logger.Close()
package logger

import (
	"io"
	"log"
	"os"
	"syscall"
)

// _LogLevel represents the level of logging.
type _LogLevel int

// Initialize the constants for the log levels.
const (
	INFO _LogLevel = iota
	ERROR
	FATAL
)

// _GetLogLevel returns the LogLevel based on the given string.
func _GetLogLevel(level string) _LogLevel {
	switch level {
	case "INFO":
		return INFO
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// _Logger represents a leveled logger.
type _Logger struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	consoleLogger *log.Logger
	logLevel      _LogLevel
}

// _NewLogger creates a new Logger instance.
func _NewLogger(logFilePath string, errorLogFilePath string, level _LogLevel) (*_Logger, error) {
	// Open log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Open error log file
	errorLogFile, err := os.OpenFile(errorLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Create loggers
	infoLogger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	errorLogger := log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	consoleLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return &_Logger{
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		consoleLogger: consoleLogger,
		logLevel:      level,
	}, nil
}

func (l *_Logger) Info(message string) {
	if l.logLevel <= INFO {
		l.infoLogger.Println(message)
		l.consoleLogger.Println(message)
	}
}

func (l *_Logger) Error(message string) {
	if l.logLevel <= ERROR {
		l.errorLogger.Println(message)
		l.infoLogger.Println("ERROR:", message)
		l.consoleLogger.Println("ERROR: ", message)
	}
}

func (l *_Logger) Fatal(message string) {
	if l.logLevel <= FATAL {
		l.errorLogger.Println(message)
		l.infoLogger.Println("FATAL:", message)
		l.consoleLogger.Fatalf("FATAL: %s", message)
	}
}

func (l *_Logger) Panic(message string) {
	l.errorLogger.Println(message)
	l.infoLogger.Println("PANIC:", message)
	l.consoleLogger.Panic(message)
	syscall.Exit(1)
}

// Close the log files.
// Since this is a singleton instance logger, this method should
// be called only once to not cause other modules to lose the connection.
// Preferably use defer to call this method in the main function.
func (l *_Logger) Close() {
	if l.infoLogger != nil {
		_ = l.infoLogger.Writer().(*os.File).Close()
	}
	if l.errorLogger != nil {
		_ = l.errorLogger.Writer().(*os.File).Close()
	}
}

func MockLogger() *_Logger {
	return &_Logger{
		infoLogger:  log.New(io.Discard, "", 0),
		errorLogger: log.New(io.Discard, "", 0),
		// consoleLogger: log.New(os.Stdout, "", log.Ldate|log.Ltime),
		consoleLogger: log.New(io.Discard, "", 0),
		logLevel:      INFO,
	}
}

// Singleton instance of the logger.
var logger *_Logger = nil

// // Initialize the logger to be used throughout the application.
// // This method should be called only once in the main function.
func InitLogger() (*_Logger, error) {
	if logger != nil {
		return logger, error(nil)
	}
	if os.Getenv("LOGGER_DISABLED") == "true" {
		logger = MockLogger()
		return logger, nil
	}
	_logger, err := _NewLogger(
		os.Getenv("LOGGER_LOG_FILE"),
		os.Getenv("LOGGER_ERROR_LOG_FILE"),
		_GetLogLevel(os.Getenv("LOGGER_LOG_LEVEL")),
	)
	logger = _logger
	return logger, err
}

// GetLogger returns the singleton instance of the logger.
// This method should be used to get the logger instance in other modules.
func GetLogger() *_Logger {
	if logger == nil {
		InitLogger()
	}
	return logger
}
