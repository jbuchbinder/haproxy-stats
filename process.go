package main

import (
	"fmt"
	"github.com/jbuchbinder/go-gmetric/gmetric"
	"log"
	"net"
	"strconv"
	"strings"
)

func ProcessStats(stats []Stats, debug bool, conn []*net.UDPConn) {
	if len(stats) < 1 {
		log.Printf("ProcessStats: no stats")
		return
	}
	for k := range stats[0] {
		if *SkipAggregates && (strings.HasSuffix(k, "_FRONTEND") || strings.HasSuffix(k, "_BACKEND")) {
			continue
		}

		log.Printf("%s:", k)

		SendStat(k+"_session_current", sumStatsInt(stats, k, "scur"), "sessions", gmetric.SLOPE_BOTH, gmetric.VALUE_INT, conn)
		SendStat(k+"_session_rate", sumStatsInt(stats, k, "rate"), "sessions", gmetric.SLOPE_BOTH, gmetric.VALUE_INT, conn)
		SendStat(k+"_bytes_in", sumStatsDouble(stats, k, "bin"), "bytes", gmetric.SLOPE_POSITIVE, gmetric.VALUE_DOUBLE, conn)
		SendStat(k+"_bytes_out", sumStatsDouble(stats, k, "bout"), "bytes", gmetric.SLOPE_POSITIVE, gmetric.VALUE_DOUBLE, conn)
	}
}

func sumStatsDouble(stats []Stats, k, v string) string {
	var total float64
	for iter := 0; iter < len(stats); iter++ {
		parsed, _ := strconv.ParseFloat(stats[iter][k][v], 64)
		total += parsed
	}
	return fmt.Sprintf("%f", total)
}

func sumStatsInt(stats []Stats, k, v string) string {
	var total int64
	for iter := 0; iter < len(stats); iter++ {
		parsed, _ := strconv.ParseInt(stats[iter][k][v], 10, 64)
		total += parsed
	}
	return fmt.Sprintf("%d", total)
}

func SendStat(name, value, units string, slope, mtype uint32, conn []*net.UDPConn) {
	if value == "" {
		return
	}
	log.Printf("SendStat %s %s %s", name, value, units)
	gm.SendMetricPackets(name, value, mtype, units, slope, uint32(*Interval*2), uint32(*Interval*2), "haproxy", gmetric.PACKET_BOTH, conn)
}
