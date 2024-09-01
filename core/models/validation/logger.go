package Sanitize

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

const (
	InfoType    = "info"
	ErrorType   = "error"
	SuccessType = "success"
	WarningType = "warning"
	ServerType  = "server"
)

// Define colors (for terminal output)
const (
	InfoColor    = "\033[34m" // Blue
	ErrorColor   = "\033[31m" // Red
	SuccessColor = "\033[32m" // Green
	WarningColor = "\033[33m" // Yellow
	ResetColor   = "\033[0m"  // Reset
)

// Create a logger
var (
	logger = log.New(os.Stdout, "", 0) // No prefix, will be added manually
)

// Custom log function
func LogMessage(logType interface{}, message string, errmsg interface{}) {
	var prefix string
	var color string

	// Convert logType to string and handle accordingly
	var logTypeStr string
	switch v := logType.(type) {
	case string:
		logTypeStr = v
	case bool:
		if v {
			logTypeStr = InfoType
		} else {
			logTypeStr = ErrorType
		}
	default:
		logTypeStr = "unknown"
	}

	// Set prefix and color based on logTypeStr
	switch logTypeStr {
	case InfoType:
		prefix = "INFO: "
		color = InfoColor
	case ErrorType:
		prefix = "ERROR: "
		color = ErrorColor
	case SuccessType:
		prefix = "SUCCESS: "
		color = SuccessColor
	case WarningType:
		prefix = "WARNING: "
		color = WarningColor
	default:
		prefix = "UNKNOWN: "
		color = ResetColor
	}

	// Get the file and line number of the caller
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// Format the filename to only include the base name and line number
	file = fmt.Sprintf("%s%d", strings.TrimPrefix(file, "/home/ubuntu/projects/websites/chatapp/"), line)

	// Handle errmsg
	switch v := errmsg.(type) {
	case error:
		if v != nil {
			message = fmt.Sprintf("%s: %v", message, v)
		}
	case string:
		if v != "" {
			message = fmt.Sprintf("%s: %s", message, v)
		}
	}

	// Format log message with file and line number at the beginning
	logMsg := fmt.Sprintf("[%s] %s%s %s%s", file, color, prefix, message, ResetColor)
	logger.Println(logMsg)
}
