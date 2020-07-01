package timex

import (
	"database/sql/driver"
	"encoding"
	"fmt"
	"time"

	"github.com/gopub/conv"
)

var (
	_ encoding.TextMarshaler   = (*Date)(nil)
	_ encoding.TextUnmarshaler = (*Date)(nil)
)

type Date struct {
	year    int
	month   int
	day     int
	weekday int

	t time.Time
}

func NewDate(year, month, day int) *Date {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return DateWithTime(t)
}

func DateWithUnix(seconds int64) *Date {
	t := time.Unix(seconds, 0)
	return DateWithTime(t)
}

func DateWithTime(t time.Time) *Date {
	return &Date{
		year:    t.Year(),
		month:   int(t.Month()),
		day:     t.Day(),
		weekday: int(t.Weekday()),
		t:       time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
	}
}

func Today() *Date {
	return DateWithTime(time.Now())
}

func Tomorrow() *Date {
	return Today().Add(0, 0, 1)
}

func Yesterday() *Date {
	return Today().Add(0, 0, -1)
}

func (d *Date) Year() int {
	return d.year
}

func (d *Date) Month() int {
	return d.month
}

func (d *Date) Day() int {
	return d.day
}

func (d *Date) Weekday() int {
	return d.weekday
}

func (d *Date) Unix() int64 {
	return d.t.Unix()
}

func (d *Date) Add(years, months, days int) *Date {
	return DateWithTime(d.t.AddDate(years, months, days))
}

func (d *Date) Time(hours, minutes int) *Time {
	return &Time{
		t: d.t.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute),
	}
}

func (d *Date) Equals(date *Date) bool {
	return d.year == date.year && d.month == date.month && d.day == date.day
}

func (d *Date) Before(date *Date) bool {
	return d.Unix() < date.Unix()
}

func (d *Date) After(date *Date) bool {
	return d.Unix() > date.Unix()
}

func (d *Date) Begin() time.Time {
	return d.t
}

func (d *Date) End() time.Time {
	return time.Date(d.year, time.Month(d.month), d.day, 23, 59, 59, 999999999, d.t.Location())
}

func (d *Date) IsToday() bool {
	return Today().Equals(d)
}

func (d *Date) IsTomorrow() bool {
	return Tomorrow().Equals(d)
}

func (d *Date) IsYesterday() bool {
	return Yesterday().Equals(d)
}

func (d *Date) String() string {
	return fmt.Sprintf("%d/%d/%d", d.year, d.month, d.day)
}

func (d *Date) PrettyText() string {
	var s string
	if IsSimplifiedChinese() {
		s = fmt.Sprintf("%d月%d日", d.month, d.day)
	} else {
		s = fmt.Sprintf("%s %d", time.Month(d.month).String()[:3], d.day)
	}
	if d.year == time.Now().Year() {
		return s
	}
	if IsSimplifiedChinese() {
		return fmt.Sprintf("%d年%s", d.year, s)
	}
	return fmt.Sprintf("%s, %d", s, d.year)
}

func (d *Date) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

func (d *Date) UnmarshalText(text []byte) error {
	s := string(text)
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = *DateWithTime(t)
	return nil
}

func (d *Date) Scan(src interface{}) error {
	t, ok := src.(time.Time)
	if !ok {
		return fmt.Errorf("expect time.Time instead of %T", src)
	}
	*d = *DateWithTime(t)
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return d.t, nil
}

func (d *Date) NextRepeat(r Repeat) *Date {
	switch r {
	case Daily:
		return d.Add(0, 0, 1)
	case Weekly:
		return d.Add(0, 0, 7)
	case Monthly:
		return d.Add(0, 1, 0)
	case Yearly:
		return d.Add(1, 0, 0)
	default:
		return nil
	}
}

func (d *Date) ShortText() string {
	if conv.AbsDuration(d.t.Sub(time.Now())) <= 2*Day {
		if IsSimplifiedChinese() {
			switch {
			case d.IsYesterday():
				return "昨天"
			case d.IsToday():
				return "今天"
			case d.IsTomorrow():
				return "明天"
			}
		} else {
			switch {
			case d.IsYesterday():
				return "Yesterday"
			case d.IsToday():
				return "Today"
			case d.IsTomorrow():
				return "Tomorrow"
			}
		}
	}
	return fmt.Sprintf("%s %s", GetWeekdaySymbol(d.weekday), d.PrettyText())
}

func (d *Date) LongText() string {
	if conv.AbsDuration(d.t.Sub(time.Now())) <= 2*Day {
		if IsSimplifiedChinese() {
			switch {
			case d.IsYesterday():
				return "昨天 " + d.PrettyText()
			case d.IsToday():
				return "今天 " + d.PrettyText()
			case d.IsTomorrow():
				return "明天 " + d.PrettyText()
			}
		} else {
			switch {
			case d.IsYesterday():
				return "Yesterday " + d.PrettyText()
			case d.IsToday():
				return "Today " + d.PrettyText()
			case d.IsTomorrow():
				return "Tomorrow " + d.PrettyText()
			}
		}
	}
	return fmt.Sprintf("%s %s", GetWeekdaySymbol(d.weekday), d.PrettyText())
}

func (d *Date) Range() *Range {
	return NewRange(d.Begin(), d.End().Add(time.Nanosecond))
}
