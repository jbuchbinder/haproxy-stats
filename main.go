package main

import (
	"flag"
	"github.com/jbuchbinder/go-gmetric/gmetric"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	SkipAggregates = flag.Bool("skipaggregates", false, "Skip FRONTEND/BACKEND values")
	StatsUrls      = flag.String("statsurl", "http://localhost:60081/", "Stats URLs or TCP addresses (CSV)")
	GmondServers   = flag.String("ganglia", "127.0.0.1:8649", "Gamglia gmond servers (host:port, CSV)")
	Debug          = flag.Bool("debug", false, "Single debugging run, no loop")
	Interval       = flag.Int("interval", 20, "Poll interval in seconds")

	gm *gmetric.Gmetric
)

func main() {
	flag.Parse()
	if *Debug {
		stats, err := GetStats()
		if err != nil {
			panic(err)
		}
		log.Printf("%v\n", stats)
		ProcessStats(stats, true, nil)
		return
	}

	log.Print("main(): Spinning up gmetric connection(s)")
	gm = &gmetric.Gmetric{}

	servers := strings.Split(*GmondServers, ",")
	for i := range servers {
		parts := strings.Split(servers[i], ":")
		port, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("Add gmetric server %s:%d", parts[0], port)
		gm.AddServer(gmetric.Server{net.ParseIP(parts[0]), int(port)})
	}

	log.Print("main(): Entering loop")
	for {
		stats, err := GetStats()
		if err != nil {
			log.Printf("Error: %v", err)
		} else {
			log.Printf("Fetched %d stats\n", len(stats))
			conn := gm.OpenConnections()
			ProcessStats(stats, false, conn)
			gm.CloseConnections(conn)
		}
		log.Printf("Going dormant for %d seconds", *Interval)
		time.Sleep(time.Duration(*Interval) * time.Second)
	}
}
