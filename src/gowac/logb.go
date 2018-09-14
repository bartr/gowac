package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var reqLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// setupLogb - initialize the log file and add multi writer
func setupLogb(logFileName string) error {
	logFileName = strings.TrimSpace(logFileName)

	if logFileName == "" {
		return errors.New("ERROR: logbpath cannot be blank")
	}

	// open the logfile
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	// setup the multi writer
	wrt := io.MultiWriter(os.Stdout, logFile)
	reqLog.SetOutput(wrt)

	return nil
}

// logHandler - http handler that writes to log file(s)
// this handler should be chained with actual request handlers
func logHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wr := &ResponseLogger{
			ResponseWriter: w,
			status:         0,
			start:          time.Now().UTC()}

		if next != nil {
			next.ServeHTTP(wr, r)
		}

		reqLog.Println(wr.status, time.Now().UTC().Sub(wr.start).Nanoseconds()/100000, r.Method, r.URL.Path, r.URL.RawQuery)
	})
}

const maxLength = 10

// ResponseLogger - wrap http.ResponseWriter to include status and size
type ResponseLogger struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

// WriteHeader - wraps http.WriteHeader
func (wr *ResponseLogger) WriteHeader(status int) {
	// store status for logging
	wr.status = status

	wr.ResponseWriter.WriteHeader(status)

	// for debug, turn caching off
	wr.Header().Add("Cache-Control", "no-cache")
}

// Write - wraps http.Write
func (wr *ResponseLogger) Write(buf []byte) (int, error) {
	n, err := wr.ResponseWriter.Write(buf)

	// store bytes written for logging
	if err == nil {
		wr.size += n
	}

	return n, err
}
