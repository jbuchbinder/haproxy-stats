package main

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type Stat map[string]string
type Stats map[string]Stat

func ReadStats(fp io.Reader) (Stats, error) {
	log.Print("ReadStats()\n")

	stats := map[string]Stat{}
	var columns []string

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		//log.Println(line)

		if strings.HasPrefix(line, "# ") {
			// Define columns
			tempcolumns := strings.Split(strings.TrimSpace(strings.TrimPrefix(line, "# ")), ",")
			for iter := 0; iter < len(tempcolumns); iter++ {
				if tempcolumns[iter] != "" {
					columns = append(columns, strings.TrimSpace(tempcolumns[iter]))
				}
			}
			log.Printf("Found %d columns : %v\n", len(columns), columns)
			continue
		}

		// Otherwise, parse values
		stat := map[string]string{}
		segments := strings.Split(line, ",")
		if len(segments) < 2 {
			continue
		}
		name := segments[0] + "_" + segments[1]
		for iter := 2; iter < len(segments)-1; iter++ {
			stat[columns[iter]] = segments[iter]
		}
		log.Printf("[%s] Found %d values : %v\n", name, len(stat), stat)

		// Pass back values
		stats[name] = stat
	}

	if err := scanner.Err(); err != nil {
		log.Printf(err.Error())
		return stats, err
	}

	return stats, nil
}
