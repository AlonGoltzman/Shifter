package logging

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"shifter.com/internal/logging_utils"
)

//====================================================
// Global variables

//====================================================
// Splunk client definition

type SplunkClient struct {
	token      string
	source     string
	sourceType string
	index      string

	logger *log.Logger
}

//Logs a message to a file according to the given information and log level.
func (client *SplunkClient) log(info map[string]interface{}, level log.Level) error {
	delete(info, "level")
	msg := info["content"]
	delete(info, "content")
	if msg == nil {
		msg = ""
	}
	switch level {
	case log.DebugLevel:
		client.logger.WithFields(log.Fields(info)).Debug(msg)
		break
	case log.WarnLevel:
		client.logger.WithFields(log.Fields(info)).Warn(msg)
		break
	case log.ErrorLevel:
		client.logger.WithFields(log.Fields(info)).Error(msg)
		break
	case log.FatalLevel:
		client.logger.WithFields(log.Fields(info)).Fatal(msg)
		break
	case log.InfoLevel:
		client.logger.WithFields(log.Fields(info)).Info(msg)
		break
	default:
		return &SplunkLoggerError{err: "Unknown logging level."}
	}


	return nil
}

//====================================================
// Errors

type SplunkLoggerError struct {
	err string
}

func (l *SplunkLoggerError) Error() string {
	return l.err
}

//====================================================
// Logger definition

type Logger struct {
	uuid    string
	logPath string

	level log.Level

	client *SplunkClient
}

// Creates a log file with a given UUID and prepares the SplunkClient to transmit information to the SplunkServer.
func OpenLogger(uuid string, source string, index string, level log.Level) (func(map[string]interface{}), error) {
	//Validations of path and uuid
	if uuid == "" {
		return nil, &SplunkLoggerError{err: "Provided empty uuid as argument."}
	}
	loggingPath := os.Getenv("SHIFTER_LOG_DIR")
	if _, err := os.Stat(loggingPath); err != nil && os.IsNotExist(err) {
		os.Mkdir(loggingPath, 0777)
	}

	var logFilePath = filepath.FromSlash(loggingPath + "/" + uuid + ".log")

	//Checking if file exists:
	var file *os.File
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	//Creating the logger object.
	l := &Logger{uuid: uuid, logPath: logFilePath, level: level}

	// Creating the splunk client.
	l.client = &SplunkClient{token: "eee30149-5dff-43b1-8eac-5fb42a645371",
		source: source,
		sourceType: "_json",
		index: index,
		logger: log.New()}

	// Configuring the logrus logger
	l.client.logger.SetLevel(level)
	l.client.logger.SetOutput(file)
	l.client.logger.Formatter = &log.JSONFormatter{}

	// Logs a test log
	l.client.log(logging_utils.NewFields().
		SetAction(logging_utils.ParseAction(logging_utils.CreateLogger)).
		SetLogType(logging_utils.ParseLogType(logging_utils.SystemLogs)).
		Finalize(),
		log.DebugLevel)

	//Returns a function to perform prints using the log function. This will invoke a goroutine
	return func(fields map[string]interface{}) {
		go l.client.log(fields, fields["level"].(log.Level))
	}, nil
}
