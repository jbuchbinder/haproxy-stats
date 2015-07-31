package main

import (
	"os"
	"testing"
)

func TestTimeConsuming(t *testing.T) {
	fp, err := os.Open("haproxy.csv")
	if err != nil {
		t.Error(err)
	}
	defer fp.Close()
	stats, err := ReadStats(fp)
	if err != nil {
		t.Error(err)
	}
	if stats["fe_production_web_http_FRONTEND"]["stot"] != "625019" {
		t.Errorf("Found %s for stot, expected 625019", stats["fe_production_web_http_FRONTEND"]["stot"])
	}
}

func test() {
}
