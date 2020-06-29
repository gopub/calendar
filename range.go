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
	begin time.Time // inclusive
	end   time.Time // exclusive
}

func NewRange(begin, end time.Time) *Range {
	begin = begin.Local()
	end = end.Local()
	r := &Range{
		begin: begin,
		end:   end,
	}
	if r.begin.After(r.end) {
		panic("timex: require begin <= end")
	}
	return r
}

func (r *Range) SetBegin(t time.Time) {
	if !t.Before(r.end) {
		panic("timex: begin must be before end")
	}
	r.begin = t
}

func (r *Range) SetEnd(t time.Time) {
	if !t.After(r.begin) {
		panic("timex: end must be after begin")
	}
	r.end = t
}

func (r *Range) Begin() time.Time {
	return r.begin
}

func (r *Range) End() time.Time {
	return r.end
}

func (r *Range) BeginUnix() int64 {
	return r.begin.Unix()
}

func (r *Range) EndUnix() int64 {
	return r.end.Unix()
}

func (r *Range) Overlap(ra *Range) bool {
	subBegin := r.begin.Sub(ra.begin)
	switch {
	case subBegin < 0: // r.begin < ra.begin
		return r.end.After(ra.begin) // r.end > ra.begin
	case subBegin > 0: // r.begin > ra.begin
		return !r.end.After(ra.end) // r.end <= ra.end
	default: // r.begin == ra.begin
		return true
	}
}

func (r *Range) Equals(ra *Range) bool {
	return r.begin.Equal(ra.begin) && r.end.Equal(ra.end)
}

func (r *Range) Contains(ra *Range) bool {
	// r.begin <= ra.begin && r.end >= ra.end
	return !r.begin.After(ra.begin) && !r.end.Before(ra.end)
}

func (r *Range) ContainsTime(t time.Time) bool {
	return !r.begin.After(t) && t.Before(r.end)
}

func (r *Range) Intersects(ra *Range) *Range {
	begin, end := r.begin, r.end
	if ra.begin.After(begin) {
		begin = ra.begin
	}
	if ra.end.Before(end) {
		end = ra.end
	}
	if begin.After(end) {
		return nil
	}
	return NewRange(begin, end)
}

func (r *Range) Duration() time.Duration {
	return r.end.Sub(r.begin)
}

func (r *Range) AddDate(years, months, days int) *Range {
	return NewRange(r.begin.AddDate(years, months, days), r.end.AddDate(years, months, days))
}

func (r *Range) IsAllDay() bool {
	return r.InDay() && r.Duration() == time.Hour*24
}

func (r *Range) InDay() bool {
	y1, m1, d1 := r.begin.Date()
	y2, m2, d2 := r.end.Add(-time.Nanosecond).Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (r *Range) Dates() []*Date {
	begin := DateWithTime(r.begin)
	end := DateWithTime(r.end)
	var l []*Date
	for d := begin; d.Equals(begin) || d.Before(end); d = d.Add(0, 0, 1) {
		l = append(l, d)
	}
	return l
}

func (r *Range) SplitInDay() []*Range {
	dates := r.Dates()
	l := make([]*Range, len(dates))
	beginTime, endTime := GetDayTime(r.begin), GetDayTime(r.end)
	for i, d := range dates {
		begin, end := d.Begin(), d.End().Add(time.Nanosecond)
		if i == 0 {
			begin = begin.Add(beginTime)
		}
		if i == len(dates)-1 {
			if endTime == 0 {
				endTime = Day
			}
			end = d.Begin().Add(endTime)
		}
		l[i] = NewRange(begin, end)
	}
	return l
}

func (r *Range) UnmarshalJSON(b []byte) error {
	var rr struct {
		Begin time.Time `json:"begin"`
		End   time.Time `json:"end"`
	}
	err := json.Unmarshal(b, &rr)
	if err != nil {
		return err
	}
	r.begin = rr.Begin
	r.end = rr.End
	return nil
}

func (r *Range) MarshalJSON() ([]byte, error) {
	var rr struct {
		Begin time.Time `json:"begin"`
		End   time.Time `json:"end"`
	}
	rr.Begin = r.begin
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
	r.begin, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[0]))
	if err != nil {
		return fmt.Errorf("parse begin %s: %w", fields[0], err)
	}
	r.end, err = time.Parse(sqlTimeLayout, strings.TrimSpace(fields[1]))
	if err != nil {
		return fmt.Errorf("parse begin %s: %w", fields[1], err)
	}
	if r.begin.After(r.end) {
		return fmt.Errorf("begin %v is after end %v", r.begin, r.end)
	}
	r.begin = r.begin.Local()
	r.end = r.end.Local()
	return nil
}

func (r Range) Value() (driver.Value, error) {
	return fmt.Sprintf("[%s, %s]", r.begin.UTC().Format(sqlTimeLayout), r.end.UTC().Format(sqlTimeLayout)), nil
}

func (r *Range) String() string {
	return fmt.Sprintf("[%s, %s)", r.begin.Format(timeLayout), r.end.Format(timeLayout))
}

func (r *Range) In(loc *time.Location) *Range {
	return &Range{
		begin: r.begin.In(loc),
		end:   r.end.In(loc),
	}
}
