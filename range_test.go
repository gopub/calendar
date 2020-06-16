package timex_test

import (
	"github.com/gopub/log"
	"github.com/gopub/timex"
	"testing"
	"time"
)

func TestRange_SplitInDay(t *testing.T) {
	tz := time.FixedZone("PST", -7*3600)
	start := time.Date(2002, 5, 3, 17, 0, 0, 0, tz)
	end := time.Date(2002, 5, 3, 18, 0, 0, 0, tz)
	r := timex.NewRange(start, end)
	for _, dr := range r.SplitInDay() {
		log.Debug(dr.Start(), dr.End())
	}
}
