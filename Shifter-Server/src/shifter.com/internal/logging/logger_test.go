package logging

import (
	"testing"
	"github.com/fuyufjh/splunk-hec-go"
	"net/http"
	"crypto/tls"
	"log"
	"time"
	"strconv"
	"github.com/sirupsen/logrus"
	"shifter.com/internal/logging_utils"
)

// A test to send a test event to Splunk.
func TestSplunkLog(t *testing.T) {
	client := hec.NewClient("https://localhost:8088", "eee30149-5dff-43b1-8eac-5fb42a645371")
	client.SetHTTPClient(&http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}})

	event1 := &hec.Event{
		Index:      String("test"),
		Source:     String("Test"),
		SourceType: String("_json"),
		Host:       String("localhost"),
		Time:       String(strconv.FormatInt(time.Now().Unix(), 10)),
		Event:      "{\"level\":\"debug\",\"msg\":\"Something failed but I'm not quitting.\"}",
	}
	err := client.WriteBatch([]*hec.Event{event1})
	if err != nil {
		log.Fatal(err)
	}
}

// A test to create a local event.
func TestSplunkLoggerLocal(t *testing.T) {
	fn, err := OpenLogger("TestUUID", "Test", "test", logrus.DebugLevel)
	if err != nil {
		t.Error(err.Error())
	}

	fn(logging_utils.NewFields().
		SetAction(logging_utils.ParseAction(logging_utils.LoggerTest)).
		SetLogType(logging_utils.ParseLogType(logging_utils.SystemTests)).
		SetLogLevel(logrus.DebugLevel).
		Finalize())

	time.Sleep(2 * time.Second)

	fn(logging_utils.NewFields().
		SetAction(logging_utils.ParseAction(logging_utils.LoggerTest)).
		SetLogType(logging_utils.ParseLogType(logging_utils.SystemTests)).
		SetLogLevel(logrus.DebugLevel).
		SetContent("This is a long test.").
		Finalize())
}

func String(str string) *string {
	return &str
}
