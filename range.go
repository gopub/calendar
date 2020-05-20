package calendar

import (
	"fmt"
	"time"
)

type Range struct {
	start time.Time
	end   time.Time
}

func NewRange(start, end time.Time) *Range {
	if !start.Before(end) {
		panic("start must be less than end")
	}
	return &Range{
		start: start,
		end:   end,
	}
}

func (r *Range) Dates() []*Date {
	start := DateWithTime(r.start)
	end := DateWithTime(r.end)
	var l []*Date
	for d := start; !d.After(end); d = d.Next() {
		l = append(l, d)
	}
	return l
}

func (r *Range) DateRanges() []*DateRange {
	dates := r.Dates()
	l := make([]*DateRange, len(dates))
	for i, d := range dates {
		start, end := time.Duration(0), EndOfDay
		if i == 0 {
			start = TimeInDay(r.start)
		}
		if i == len(dates)-1 {
			end = TimeInDay(r.end)
		}
		l[i] = NewDateRange(d, start, end)
	}
	return l
}

type DateRange struct {
	date  *Date
	start time.Duration
	end   time.Duration
}

func NewDateRange(date *Date, start, end time.Duration) *DateRange {
	start = start.Round(time.Minute)
	end = end.Round(time.Minute)
	if start < 0 || start > EndOfDay {
		panic("start must be in [0, 24h)")
	}

	if end < time.Minute || end > EndOfDay+time.Nanosecond {
		panic("end must be 0 or in [1m, 24h]: " + fmt.Sprint(end))
	}

	if end-start < time.Minute {
		panic("expect: end - start >= 1m")
	}

	return &DateRange{
		date:  date,
		start: start,
		end:   end,
	}
}

func (r *DateRange) Date() *Date {
	return r.date
}

func (r *DateRange) Start() (hour, minute int) {
	return int(r.start.Hours()), int(r.start.Minutes()) % 60
}

func (r *DateRange) End() (hour, minute int) {
	return int(r.end.Hours()), int(r.end.Minutes()) % 60
}

func (r *DateRange) Duration() time.Duration {
	return r.end - r.start + time.Minute
}

func (r *DateRange) String() string {
	sh, sm := r.Start()
	eh, em := r.End()
	return fmt.Sprintf("%s %02d:%02d-%02d:%02d", r.date, sh, sm, eh, em)
}

func (r *DateRange) IsAllDay() bool {
	return r.end-r.start == time.Hour*24
}
