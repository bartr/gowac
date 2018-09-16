package logb

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TODO - should we replace with Gorilla logger?
// TODO - add Apache log file support

var reqLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// SetLogFile - initialize the log file and add multi writer
func SetLogFile(logFile string) error {
	logFile = strings.TrimSpace(logFile)

	if logFile == "" {
		return errors.New("ERROR: logbpath cannot be blank")
	}

	// open the logfile
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	// setup the multi writer
	wrt := io.MultiWriter(os.Stdout, f)
	reqLog.SetOutput(wrt)

	return nil
}

//Handler - http handler that writes to log file(s)
func Handler(next http.Handler) http.Handler {
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

// ResponseLogger - wrap http.ResponseWriter to include status and size
type ResponseLogger struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

// WriteHeader - wraps http.WriteHeader
func (r *ResponseLogger) WriteHeader(status int) {
	// store status for logging
	r.status = status

	r.ResponseWriter.WriteHeader(status)
}

// Write - wraps http.Write
func (r *ResponseLogger) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)

	// store bytes written for logging
	if err == nil {
		r.size += n
	}

	return n, err
}
