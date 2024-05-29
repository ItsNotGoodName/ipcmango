package main

import (
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/k0kubun/pp/v3"
)

func main() {
	now := time.Now()
	start, end := now.Add(-170*24*time.Hour), now
	pp.Println("DESCENDING", "Start", start, "End", end)
	r := dahua.NewFileScanRange(start, end, dahua.MaxFileScanPeriod, false)
	for r.Next() {
		start, end := r.Range()
		pp.Println(start, end, r.Percent())
	}

	pp.Println("ASCENDING", "Start", start, "End", end)
	r = dahua.NewFileScanRange(start, end, dahua.MaxFileScanPeriod, true)
	for r.Next() {
		start, end := r.Range()
		pp.Println(start, end, r.Percent())
	}
}
