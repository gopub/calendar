package timex

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gopub/conv"
)

var (
	_ encoding.TextMarshaler   = (*Range)(nil)
	_ encoding.TextUnmarshaler = (*Range)(nil)
)

type Range struct {
	start time.Time
	end   time.Time
}

func NewRange(start, end time.Time) *Range {
	r := new(Range)
	r.Set(start, end)
	return r
}

func (r *Range) Set(start, end time.Time) {
	if !start.Before(end) {
		panic("start must be less than end")
	}
	r.start = start
	r.end = end
}

func (r *Range) Start() time.Time {
	return r.start
}

func (r *Range) End() time.Time {
	return r.end
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

func (r *Range) DailyRanges() []*DailyRange {
	dates := r.Dates()
	l := make([]*DailyRange, len(dates))
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

func (r *Range) MarshalText() (text []byte, err error) {
	var rr struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	rr.Start = r.start
	rr.End = r.end
	return json.Marshal(rr)
}

func (r *Range) UnmarshalText(text []byte) error {
	var rr struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	}
	err := json.Unmarshal(text, &rr)
	if err != nil {
		return err
	}
	r.start = rr.Start
	r.end = rr.End
	return nil
}

var (
	_ driver.Valuer = (*Range)(nil)
	_ sql.Scanner   = (*Range)(nil)
)

const sqlTimeLayout = "2006-01-02 15:04:05.999999999-07"

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
	return nil
}

func (r Range) Value() (driver.Value, error) {
	return fmt.Sprintf("[%s, %s]", r.start.Format(sqlTimeLayout), r.end.Format(sqlTimeLayout)), nil
}

type DailyRange struct {
	date  *Date
	start time.Duration
	end   time.Duration
}

func NewDateRange(date *Date, start, end time.Duration) *DailyRange {
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

	return &DailyRange{
		date:  date,
		start: start,
		end:   end,
	}
}

func (r *DailyRange) Date() *Date {
	return r.date
}

func (r *DailyRange) Start() (hour, minute int) {
	return int(r.start.Hours()), int(r.start.Minutes()) % 60
}

func (r *DailyRange) End() (hour, minute int) {
	return int(r.end.Hours()), int(r.end.Minutes()) % 60
}

func (r *DailyRange) Duration() time.Duration {
	return r.end - r.start + time.Minute
}

func (r *DailyRange) IsAllDay() bool {
	return r.end-r.start == time.Hour*24
}

func (r *DailyRange) StartsBefore(dr *DailyRange) bool {
	if r.date.Before(dr.date) {
		return true
	}

	if r.date.After(dr.date) {
		return false
	}

	return r.start < dr.start
}

func (r *DailyRange) StartsAfter(dr *DailyRange) bool {
	return dr.StartsBefore(r)
}

func (r *DailyRange) EndsBefore(dr *DailyRange) bool {
	if r.date.Before(dr.date) {
		return true
	}

	if r.date.After(dr.date) {
		return false
	}

	return r.end < dr.end
}

func (r *DailyRange) EndsAfter(dr *DailyRange) bool {
	return dr.EndsBefore(r)
}

func (r *DailyRange) String() string {
	sh, sm := r.Start()
	eh, em := r.End()
	return fmt.Sprintf("%s %02d:%02d-%02d:%02d", r.date, sh, sm, eh, em)
}
