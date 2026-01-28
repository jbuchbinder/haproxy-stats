package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/jbuchbinder/go-gmetric/gmetric"
)

var (
	skipAggregates = flag.Bool("skipaggregates", false, "Skip FRONTEND/BACKEND values")
	statsUrls      = flag.String("statsurl", "http://localhost:60081/", "Stats URLs or TCP addresses (CSV)")
	gmondServers   = flag.String("ganglia", "127.0.0.1:8649", "Gamglia gmond servers (host:port, CSV)")
	debug          = flag.Bool("debug", false, "Single debugging run, no loop")
	interval       = flag.Int("interval", 20, "Poll interval in seconds")
	timeout        = flag.Duration("timeout", 5*time.Second, "Connection timeout")
	daemonize      = flag.Bool("daemon", false, "Daemonize")

	gm *gmetric.Gmetric
)

func main() {
	flag.Parse()

	// Daemon stuff if we're configured for it.
	if *daemonize {
		go func() {
			log.Printf("main(): INFO: Spawning systemd integration")

			interval, err := daemon.SdWatchdogEnabled(false)
			if err != nil {
				log.Printf("ERR: %s", err.Error())
				return
			}
			if interval == 0 {
				log.Printf("ERR: interval == 0")
				return
			}
			for {
				daemon.SdNotify(false, daemon.SdNotifyWatchdog)
				time.Sleep(interval / 3)
			}
		}()
	}

	if *debug {
		stats, err := getStats()
		if err != nil {
			panic(err)
		}
		log.Printf("%v\n", stats)
		processStats(stats, true, nil)
		return
	}

	log.Print("main(): INFO: Spinning up gmetric connection(s)")
	gm = &gmetric.Gmetric{}

	servers := strings.Split(*gmondServers, ",")
	for i := range servers {
		parts := strings.Split(servers[i], ":")
		port, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("Add gmetric server %s:%d", parts[0], port)
		gm.AddServer(gmetric.Server{Server: net.ParseIP(parts[0]), Port: int(port)})
	}

	log.Print("main(): INFO: Entering loop")
	for {
		stats, err := getStats()
		if err != nil {
			log.Printf("main(): ERR: %v", err)
		} else {
			log.Printf("Fetched %d stats\n", len(stats))
			conn := gm.OpenConnections()
			processStats(stats, false, conn)
			gm.CloseConnections(conn)
		}
		log.Printf("Going dormant for %d seconds", *interval)
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
