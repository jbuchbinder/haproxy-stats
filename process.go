package main

import (
	"fmt"
	"github.com/jbuchbinder/go-gmetric/gmetric"
	"log"
	"net"
	"strconv"
	"strings"
)

func processStats(mystats []stats, debug bool, conn []*net.UDPConn) {
	if len(mystats) < 1 {
		log.Printf("processStats: no stats")
		return
	}
	log.Printf("processStats: Aggregating from %d stats sources", len(mystats))
	for k := range mystats[0] {
		if *skipAggregates && (strings.HasSuffix(k, "_FRONTEND") || strings.HasSuffix(k, "_BACKEND")) {
			continue
		}

		log.Printf("%s:", k)

		sendStat(k+"_session_current", sumStatsInt(mystats, k, "scur"), "sessions", gmetric.SLOPE_BOTH, gmetric.VALUE_INT, conn)
		sendStat(k+"_session_rate", sumStatsInt(mystats, k, "rate"), "sessions", gmetric.SLOPE_BOTH, gmetric.VALUE_INT, conn)
		sendStat(k+"_bytes_in", sumStatsDouble(mystats, k, "bin"), "bytes", gmetric.SLOPE_POSITIVE, gmetric.VALUE_DOUBLE, conn)
		sendStat(k+"_bytes_out", sumStatsDouble(mystats, k, "bout"), "bytes", gmetric.SLOPE_POSITIVE, gmetric.VALUE_DOUBLE, conn)
	}
}

func sumStatsDouble(mystats []stats, k, v string) string {
	var total float64
	for iter := 0; iter < len(mystats); iter++ {
		parsed, _ := strconv.ParseFloat(mystats[iter][k][v], 64)
		total += parsed
	}
	return fmt.Sprintf("%f", total)
}

func sumStatsInt(mystats []stats, k, v string) string {
	var total int64
	for iter := 0; iter < len(mystats); iter++ {
		parsed, _ := strconv.ParseInt(mystats[iter][k][v], 10, 64)
		total += parsed
	}
	return fmt.Sprintf("%d", total)
}

func sendStat(name, value, units string, slope, mtype uint32, conn []*net.UDPConn) {
	if value == "" {
		return
	}
	log.Printf("sendStat %s %s %s", name, value, units)
	gm.SendMetricPackets(name, value, mtype, units, slope, uint32(*interval*2), uint32(*interval*2), "haproxy", gmetric.PACKET_BOTH, conn)
}
