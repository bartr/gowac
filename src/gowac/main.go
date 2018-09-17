package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bartr/gohandlers/logb"
)

var logPath = "/home/LogFiles/"
var port = 8080

func main() {

	parseCommandLine()

	// setup handlers
	http.HandleFunc("/requests", httpDumpRequests)
	http.Handle("/", logb.Handler(rawRequestHandler(http.HandlerFunc(httpDefault))))

	// setup log files
	if err := setupAppLog(buildFullLogName(logPath, "app", ".log")); err != nil {
		log.Fatal(err)
	}

	if err := logb.SetLogFile(buildFullLogName(logPath, "request", ".log")); err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on", port)
	log.Println("Logging to", logPath)

	// run the web server
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}

// parseCommandLine -port and -logpath are supported
func parseCommandLine() {
	// /home/LogFiles is the shared CIFS directory on App Services
	// log in current directory if not running in App Services
	if _, err := os.Stat(logPath); err != nil {
		logPath = "./"
	}

	// parse flags
	lfp := flag.String("logpath", "", "path to log files")
	p := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	// set log path
	if *lfp != "" {
		logPath = *lfp
	}

	// set port
	if *p <= 0 || *p >= 64*1024 {
		flag.Usage()
		log.Fatal("invalid port")
	}
	port = *p
}

// setup log multi writer
func setupAppLog(logFile string) error {

	// prepend date and time to log entries
	log.SetFlags(log.Ldate | log.Ltime)

	// open the log file
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	// setup a multiwriter to log to file and stdout
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	return nil
}

// build the full log file name
// app services sets the WEBSITE_ROLE_INSTANCE_ID environment variable
// since we're writing to the CIFS share, we need to differentiate log file names
//   in case there are multiple instances running
func buildFullLogName(logPath string, logPrefix string, logExtension string) string {
	if !strings.HasSuffix(logPath, "/") {
		logPath += "/"
	}

	fileName := logPath + logPrefix

	// use instance ID to differentiate log files between instances in App Services
	if iid := os.Getenv("WEBSITE_ROLE_INSTANCE_ID"); iid != "" {
		fileName += "_" + strings.TrimSpace(iid)
	}

	// make sure logExtension has a period
	if !strings.HasPrefix(logExtension, ".") {
		logExtension = "." + logExtension
	}

	return fileName + logExtension
}
