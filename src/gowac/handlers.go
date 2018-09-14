package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// handle all requests
func httpDefault(w http.ResponseWriter, r *http.Request) {

	s := strings.ToLower(r.URL.Path)

	// handle default web page
	if s == "/" || strings.HasPrefix(s, "/index.") || strings.HasPrefix(s, "/default.") {
		http.ServeFile(w, r, "www/default.html")
		return
	}

	// handle /home/LogFiles browsing
	if strings.HasPrefix(s, "/home") && strings.HasPrefix(logPath, "/home/") {
		http.ServeFile(w, r, r.URL.Path)
		return
	}

	// don't allow directory browsing (unless you want to)
	if strings.HasSuffix(s, "/") {
		w.WriteHeader(403)
		return
	}

	// serve the file from the www directory
	http.ServeFile(w, r, "www"+s)
}

// use http dump request to show the full request
func httpDumpRequests(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintln(requests)

	// trim the []
	s = s[1 : len(s)-2]

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(s))
}

var requests []string

// rawRequestHandler - http handler that saves the raw request
// this handler can be chained with other handlers
func rawRequestHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s := strings.ToLower(r.URL.Path)

		if !strings.HasSuffix(s, "favicon.ico") {

			// get the full request
			b, err := httputil.DumpRequest(r, true)

			if err != nil {
				log.Printf("Error: %s\n", b)
				return
			}

			s := string(b) + "====================\n"

			// prepend to array
			arr := make([]string, len(requests)+1)
			arr[0] = s
			copy(arr[1:], requests)
			requests = arr

			// slice array to maxLength
			if len(requests) > maxLength {
				requests = requests[:maxLength]
			}
		}

		// call the next handler
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
