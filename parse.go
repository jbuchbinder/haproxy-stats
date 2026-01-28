package main

import (
	"bufio"
	"io"
	"log"
	"strings"
)

type stat map[string]string
type stats map[string]stat

func readStats(fp io.Reader) (stats, error) {
	log.Print("readStats()\n")

	mystats := map[string]stat{}
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
		mystat := map[string]string{}
		segments := strings.Split(line, ",")
		if len(segments) < 2 {
			continue
		}
		name := segments[0] + "_" + segments[1]
		for iter := 2; iter < len(segments)-1; iter++ {
			mystat[columns[iter]] = segments[iter]
		}
		log.Printf("[%s] Found %d values : %v\n", name, len(mystat), mystat)

		// Pass back values
		mystats[name] = mystat
	}

	if err := scanner.Err(); err != nil {
		log.Printf("ERR: %s", err.Error())
		return mystats, err
	}

	return mystats, nil
}
