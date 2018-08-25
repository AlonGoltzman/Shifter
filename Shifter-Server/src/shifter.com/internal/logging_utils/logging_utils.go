package logging_utils

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
	"fmt"
)

type Fields struct {
	fields map[string]interface{}
}

//Generates an empty fields set for the user to manipulate.
func NewFields() *Fields {
	return &Fields{fields: map[string]interface{}{}}
}

//Sets the action type, whether it be a system operation action like creating a logger or a user action like submitting shifts.
func (fields *Fields) SetAction(action *Action) *Fields {
	fields.fields["action"] = action.String()
	return fields
}

//Sets the log type, helps filter system and user actions.
func (fields *Fields) SetLogType(logType *LogType) *Fields {
	fields.fields["log_type"] = logType.String()
	return fields
}

//Sets the log level, will be used by the SplunkClient in "splunk_logger" in order to print in the correct level.
func (fields *Fields) SetLogLevel(level log.Level) *Fields {
	fields.fields["level"] = level
	return fields
}

//Sets the 'msg' tag of the JSON for further information if needed.
func (fields *Fields) SetContent(msg string) *Fields {
	fields.fields["content"] = msg
	return fields
}

//Finalizes the fields will adding an epoch for time management.
func (fields *Fields) Finalize() map[string]interface{} {
	fields.fields["epoch"] = strconv.FormatInt(time.Now().Unix(), 10)
	return fields.fields
}

//====================================================
// Actions

type Action struct {
	value  int
	action string
}

const (
	CreateLogger = iota
	LoggerTest
)

func ParseAction(val int) *Action {
	switch val {
	case CreateLogger:
		return &Action{value: val, action: "Create Logger"}
	case LoggerTest:
		return &Action{value: val, action: "Testing Logger"}
	}

	return &Action{value: 0, action: fmt.Sprintf("Unknown [%v]", val)}
}

func (act *Action) String() string {
	return act.action
}

//====================================================
// Log Types

type LogType struct {
	value   int
	logType string
}

const (
	SystemLogs  = iota
	SystemTests
)

func ParseLogType(val int) *LogType {
	switch val {
	case SystemLogs:
		return &LogType{value: val, logType: "System"}
	case SystemTests:
		return &LogType{value: val, logType: "SysTests"}
	}

	return &LogType{value: 0, logType: fmt.Sprintf("Unknown [%v]", val)}
}

func (logType *LogType) String() string {
	return logType.logType
}
