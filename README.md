# HAPROXY-STATS

[![Status](https://secure.travis-ci.org/jbuchbinder/haproxy-stats.png)](http://travis-ci.org/jbuchbinder/haproxy-stats) [![Build Status](https://drone.io/github.com/jbuchbinder/haproxy-stats/status.png)](https://drone.io/github.com/jbuchbinder/haproxy-stats/latest)

[haproxy](http://www.haproxy.org/) load balancer stats collector for Ganglia. Supports multiple stats sockets and Ganglia gmond collectors.

## BUILDING

	go get -d
	go build

## USAGE

```
Usage of ./haproxy-stats:
  -debug
    	Single debugging run, no loop
  -ganglia string
    	Gamglia gmond servers (host:port, CSV) (default "127.0.0.1:8649")
  -interval int
    	Poll interval in seconds (default 20)
  -skipaggregates
    	Skip FRONTEND/BACKEND values
  -statsurl string
    	Stats URLs or TCP addresses (CSV) (default "http://localhost:60081/")
```

