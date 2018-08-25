package logging
//
//import (
//	"time"
//	"fmt"
//	"os"
//	"sync"
//	"strings"
//	"strconv"
//	"path/filepath"
//	"io/ioutil"
//)
//
////====================================================
//// Global variables
//var LevelList []*Level = []*Level{
//	&Level{value: 0, name: "OFF"},
//	&Level{value: 1, name: "FATAL"},
//	&Level{value: 2, name: "ERROR"},
//	&Level{value: 3, name: "WARN"},
//	&Level{value: 4, name: "INFO"},
//	&Level{value: 5, name: "DEBUG"},
//	&Level{value: 6, name: "TRACE"},
//}
//
//var dateFormat = time.StampMilli
//var defaultLogLevel = LevelList[4]
//var loggingPath = os.Getenv("SHIFTER_LOG_DIR")
////====================================================
//// Logger error
//
//type LoggerError struct {
//	err string
//}
//
//func (loggerErr *LoggerError) Error() string {
//	return loggerErr.err
//}
//
////====================================================
//// Log Level definition
//
//type Level struct {
//	value uint8
//	name  string
//}
//
//func (lvl *Level) String() string {
//	return fmt.Sprintf("Defined Level %v[%v]", lvl.name, lvl.value)
//}
//
//func (lvl *Level) Literal() string {
//	return lvl.name
//}
//
////====================================================
//// Log structure definition
//
//type Log struct {
//	logTimestamp time.Time
//	from         string
//	lvl          Level
//	content      string
//}
//
//func (log *Log) String() string {
//	return fmt.Sprintf("[%v | %v | %v] %v", log.lvl.Literal(), log.logTimestamp.Format(dateFormat), log.from, log.content)
//}
//
////====================================================
//// Logger definition
//
//type Logger struct {
//	definedLvl  *Level
//	logFilepath string
//
//	//Sync operations in the logger
//	waitGroup sync.WaitGroup
//	input     chan *Log
//	inputBKP  chan *Log
//	inputCtrl chan struct{}
//}
//
//// The go-routine thread method.
//func (logger *Logger) run() {
//	defer logger.waitGroup.Done() // Run this after the run function is done.
//
//	for {
//		select {
//		case logMsg := <-logger.input: //Received log message.
//			ioutil.WriteFile(logger.logFilepath, []byte(logMsg.String()), 0600)
//			fmt.Println("Hi")
//			//At this point the routine is alseep.
//		case _, ok := <-logger.inputCtrl:
//			if ok {
//				continue
//			}
//			return
//		}
//	}
//}
//
//// Recursively continues to add 1 to the end of the file name (1,2,3,4,etc...) until
//// a file is found that does not already have a lock on it.
//// Warning, might cause a stack over-flow issue and propagate to the logger go-routine.
//func (logger *Logger) nextLogFile() {
//	file := logger.logFilepath
//	currentIteration, err := strconv.Atoi(strings.Split(file, ".")[len(strings.Split(file, "."))-1])
//	if err != nil {
//		file += ".1"
//	}
//	file += strconv.Itoa(currentIteration + 1)
//
//	//logger.logLock = flock.NewFlock(file)
//	//locked, err := logger.logLock.TryLock()
//	if err != nil {
//		logger.nextLogFile()
//	}
//
//	//if locked {
//	//	logger.logLock.Unlock()
//	//}
//}
//
//// Creates a logger with the given parameters or returns an error.
//func genLogger(lvl *Level, logFilePath string) (*Logger, error) {
//
//	if lvl.value == 0 { //Logger turned off
//		return nil, &LoggerError{err: fmt.Sprintf("Logger provided with \"OFF\" level.")}
//	}
//
//	if logFilePath == "" { //Log file is undefined
//		return nil, &LoggerError{err: fmt.Sprintf("Logger provided with empty filepath parameter.")}
//	}
//
//	logger := &Logger{definedLvl: lvl, logFilepath: logFilePath}
//	return logger, nil
//}
//
//// Sends a log message through to the logger
//func (logger *Logger) LogMessage(log *Log) {
//	if logger.input == nil {
//		logger.input = logger.inputBKP
//	}
//	logger.inputCtrl <- struct{}{}
//	logger.input <- log
//}
//
//// Creates a logger and returns a function or an error according to the process.
//func OpenLogger(lvl *Level, value string) (func(*Log), error) {
//	if lvl == nil {
//		lvl = defaultLogLevel
//		//TODO: Add trace "default log level set".
//	}
//
//	// Create log directory if needed:
//	if _, err := os.Stat(loggingPath); err != nil && os.IsNotExist(err) {
//		os.Mkdir(loggingPath, 0777)
//	}
//
//	var logFilePath = filepath.FromSlash(loggingPath + "/" + value + ".log")
//	//Checking if file exists:
//	_, err := os.Stat(logFilePath)
//	if os.IsNotExist(err) {
//		_, err = os.Create(logFilePath)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	//Configure logger
//	logger, err := genLogger(lvl, logFilePath)
//	if err != nil {
//		return nil, err
//	}
//	channel := make(chan *Log)
//	logger.input = channel
//	logger.inputBKP = channel
//	logger.inputCtrl = make(chan struct{})
//
//	logger.waitGroup = sync.WaitGroup{}
//	logger.waitGroup.Add(1)
//
//	go logger.run()
//	return logger.LogMessage, nil
//}
