package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bartr/gohandlers/logb"
	"github.com/bartr/gohandlers/rawrequest"
)

// default values
// this is the App Service shared CIFS folder
var logPath = "/home/LogFiles/"
var port = 8080

func main() {

	parseCommandLine()

	// this sets the log files for app and requests
	if err := setupLogging(logPath); err != nil {
		log.Fatal(err)
	}

	// setup handlers
	http.HandleFunc("/requests", rawrequest.DisplayRawRequests)
	http.Handle("/", logb.Handler(rawrequest.Handler(http.HandlerFunc(rootHandler))))

	log.Println("Listening on", port)
	log.Println("Logging to", logPath)

	// run the web server
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}

// parseCommandLine -port and -logpath are supported
func parseCommandLine() {
	// parse flags
	lfp := flag.String("logpath", "", "path to log files")
	p := flag.Int("port", port, "port to listen on")
	flag.Parse()

	// set log path
	if *lfp != "" {
		logPath = *lfp
	}

	// /home/LogFiles is the shared CIFS directory on App Services
	// log in current directory if not running in App Services
	// TODO - do we really want to do this?
	if _, err := os.Stat(logPath); err != nil {
		logPath = "./"
	}

	// set port
	if *p <= 0 || *p >= 64*1024 {
		flag.Usage()
		log.Fatal("invalid port")
	}

	port = *p
}

// setupLogging - setup the app and request log files
func setupLogging(logPath string) error {

	// setup app log
	if err := logb.SetupAppLog(logPath, "app", ".log"); err != nil {
		return err
	}

	// setup request log
	if err := logb.SetLogFile(logPath, "request", ".log"); err != nil {
		return err
	}

	return nil
}

// handle all requests
func rootHandler(w http.ResponseWriter, r *http.Request) {

	s := strings.ToLower(r.URL.Path)

	// handle default web page
	if s == "/" || strings.HasPrefix(s, "/index.") || strings.HasPrefix(s, "/default.") {
		http.ServeFile(w, r, "www/default.html")
		w.Header().Add("Cache-Control", "no-cache")
		return
	}

	// handle /home/LogFiles browsing
	if strings.HasPrefix(s, "/home") && strings.HasPrefix(logPath, "/home/") {
		// TODO - don't allow /../ in path
		http.ServeFile(w, r, r.URL.Path)
		w.Header().Add("Cache-Control", "no-cache")
		return
	}

	// don't allow directory browsing (unless you want to)
	if strings.HasSuffix(s, "/") {
		w.WriteHeader(403)
		return
	}

	// TODO - don't allow /../ in path

	// serve the file from the www directory
	http.ServeFile(w, r, "www"+s)
	w.Header().Add("Cache-Control", "no-cache")
}
