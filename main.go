package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// client is shared by every check so it's built once.
// The 5s timeout stops a slow site from hanging the program.
var client = http.Client{
	Timeout: 5 * time.Second,
}

// checkURL requests site and returns its status code, the response
// time in milliseconds, and any error.
func checkURL(site string) (int, int64, error) {

	startTime := time.Now()

	response, err := client.Get(site)

	// Time the request in milliseconds.
	timeTaken := time.Since(startTime).Milliseconds()

	// On error response is nil, so return before touching it.
	if err != nil {
		return 0, 0, err
	}

	// Close the body when we return, or the connection stays open.
	defer response.Body.Close()

	return response.StatusCode, timeTaken, nil
}

// errorLabel turns a request error into a short label.
func errorLabel(err error) string {

	var timeoutErr net.Error
	if errors.As(err, &timeoutErr) && timeoutErr.Timeout() {
		return "timeout"
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return "unreachable"
	}

	return "error: " + err.Error()
}

// main reads URLs from urls.txt, one per line, and checks each one.
func main() {
	file, err := os.Open("urls.txt")

	if err != nil {
		log.Fatal(err)
	}

	// Close the file when main returns.
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		site := strings.TrimSpace(scanner.Text())
		// Skip blank lines.
		if site != "" {
			statusCode, timeTaken, err := checkURL(site)
			// A failed site prints its error and the loop carries on.
			if err != nil {
				errLabel := errorLabel(err)
				fmt.Println(site, "-", errLabel)
			} else {
				fmt.Println(site, "-", statusCode, "-", timeTaken)
			}
		}
	}

	// Scan stops at the end of the file or on a read error; this tells them apart.
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
