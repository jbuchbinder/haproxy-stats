package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func getStats() ([]stats, error) {
	mystats := make([]stats, 0)

	urls := strings.Split(*statsUrls, ",")

	for iter := range urls {
		s, e := readSingleStats(urls[iter])
		if e == nil {
			mystats = append(mystats, s)
		}
	}

	return mystats, nil
}

func readSingleStats(url string) (stats, error) {
	var mystats stats

	// For http urls, use standard GET requests
	if strings.HasPrefix(url, "http") {
		resp, err := http.Get(url + ";csv")
		if err != nil {
			return mystats, err
		}
		defer resp.Body.Close()
		return readStats(resp.Body)
	}

	// Use TCP collection instead
	d := net.Dialer{Timeout: *timeout}
	conn, err := d.Dial("tcp", url)
	if err != nil {
		return mystats, err
	}
	defer func() {
		fmt.Fprintf(conn, "quit\n")
		conn.Close()
	}()
	_, err = fmt.Fprintf(conn, "show stat\n")
	return readStats(conn)
}
