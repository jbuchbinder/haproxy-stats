package main

import (
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
	resp, err := http.Get(url + ";csv")
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()
	return ReadStats(resp.Body)
}
