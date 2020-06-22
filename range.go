package timex

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gopub/conv"
)

var (
	_ json.Marshaler   = (*Range)(nil)
	_ json.Unmarshaler = (*Range)(nil)
)

type Range struct {
	start time.Time // inclusive
	end   time.Time // exclusive
}

func NewRange(start, end time.Time) *Range {
	start = start.Local()
	end = end.Local()
	r := &Range{
		start: start,
		end:   end,
	}
	if !r.start.Before(r.end) {
		panic("timex: start must be before end")
	}
	return r
}

func (r *Range) SetStart(t time.Time) {
	if !t.Before(r.end) {
		panic("timex: start must be before end")
	}
	r.start = t
}

func (r *Range) SetEnd(t time.Time) {
	if !t.After(r.start) {
		panic("timex: end must be after start")
	}
	r.end = t
}

func (r *Range) Start() time.Time {
	return r.start
}

func (r *Range) End() time.Time {
	return r.end
}

func (r *Range) StartsBefore(ra *Range) bool {
	return r.start.Before(ra.start)
}

func (r *Range) EndsAfter(ra *Range) bool {
	return r.end.After(ra.end)
}

func (r *Range) Includes(t time.Time) bool {
	return !r.start.After(t) && t.Before(r.end)
}

func (r *Range) Duration() time.Duration {
	return r.end.Sub(r.start)
}

func (r *Range) AddDate(years, months, days int) *Range {
	return NewRange(r.start.AddDate(years, months, days), r.end.AddDate(years, months, days))
}

func (r *Range) IsAllDay() bool {
	return r.InDay() && r.Duration() == time.Hour*24
}

func (r *Range) InDay() bool {
	y1, m1, d1 := r.start.Date()
	y2, m2, d2 := r.end.Add(-time.Nanosecond).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (r *Range) Dates() []*Date {
	start := DateWithTime(r.start)
	end := DateWithTime(r.end)
	var l []*Date
	for d := start; !d.After(end); d = d.Add(0, 0, 1) {
		l = append(l, d)
	}
	return l
}

func (r *Range) SplitInDay() []*Range {
	dates := r.Dates()
	l := make([]*Range, len(dates))
	startTime, endTime := GetDayTime(r.start), GetDayTime(r.end)
	for i, d := range dates {
		start, end := d.Start(), d.End().Add(time.Nanosecond)
		if i == 0 {
			start = start.Add(startTime)
		}
		if i == len(dates)-1 {
			if endTime == 0 {
				endTime = Day
			}
			end = d.Start().Add(endTime)
		}
		l[i] = NewRange(start, end)
	}
	return l
}

func (r *Range) UnmarshalJSON(b []byte) error {
	var rr struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	err := json.Unmarshal(b, &rr)
	if err != nil {
		return err
	}
	r.start = rr.Start
	r.end = rr.End
	return nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	var rr struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	rr.Start = r.start
	rr.End = r.end
	return json.Marshal(rr)
}

var (
	_ driver.Valuer = (*Range)(nil)
	_ sql.Scanner   = (*Range)(nil)
)

const (
	sqlTimeLayout = "2006-01-02 15:04:05.999999999-07"
	timeLayout    = "2006-01-02 15:04:05-07"
)

func (r *Range) Scan(src interface{}) error {
	s, err := conv.ToString(src)
	if err != nil {
		return err
	}

	if s == "" {
		return nil
	}

	s = strings.Replace(s, `"`, "", -1)

	if s[0] != '[' {
		return fmt.Errorf("cannot parse %s", s)
	}

	if c := s[len(s)-1]; c != ']' {
		return fmt.Errorf("cannot parse %s", s)
	}

	s = s[1 : len(s)-1]

	fields := strings.Split(s, ",")
	if len(fields) != 2 {
		return fmt.Errorf("parse composite fields %s", s)
	}
	r.start, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[0]))
	if err != nil {
		return fmt.Errorf("parse start %s: %w", fields[0], err)
	}
	r.end, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[1]))
	if err != nil {
		return fmt.Errorf("parse start %s: %w", fields[1], err)
	}
	if r.start.After(r.end) {
		return fmt.Errorf("start %v is after end %v", r.start, r.end)
	}
	r.start = r.start.Local()
	r.end = r.end.Local()
	return nil
}

func (r Range) Value() (driver.Value, error) {
	return fmt.Sprintf("[%s, %s]", r.start.UTC().Format(sqlTimeLayout), r.end.UTC().Format(sqlTimeLayout)), nil
}

func (r *Range) String() string {
	return fmt.Sprintf("[%s, %s)", r.start.Format(timeLayout), r.end.Format(timeLayout))
}

func (r *Range) In(loc *time.Location) *Range {
	return &Range{
		start: r.start.In(loc),
		end:   r.end.In(loc),
	}
}
