package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetStats() ([]Stats, error) {
	stats := make([]Stats, 0)

	urls := strings.Split(*StatsUrls, ",")

	for iter := 0; iter < len(urls); iter++ {
		s, e := ReadSingleStats(urls[iter])
		if e == nil {
			stats = append(stats, s)
		}
	}

	return stats, nil
}

func ReadSingleStats(url string) (Stats, error) {
	var stats Stats

	// For http urls, use standard GET requests
	if strings.HasPrefix(url, "http") {
		resp, err := http.Get(url + ";csv")
		if err != nil {
			return stats, err
		}
		defer resp.Body.Close()
		return ReadStats(resp.Body)
	}

	// Use TCP collection instead
	conn, err := net.Dial("tcp", url)
	if err != nil {
		return stats, err
	}
	defer func() {
		fmt.Fprintf(conn, "quit\n")
		conn.Close()
	}()
	_, err = fmt.Fprintf(conn, "show stat\n")
	return ReadStats(conn)
}
