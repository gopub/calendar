package timex

import (
	"fmt"
	"time"
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
		t:       t,
	}
}

func Today() *Date {
	return DateWithTime(time.Now())
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

func (d *Date) Next() *Date {
	return DateWithTime(d.t.Add(time.Hour * 24))
}

func (d *Date) Prev() *Date {
	return DateWithTime(d.t.Add(-time.Hour * 24))
}

func (d *Date) NextWeek() *Date {
	return DateWithTime(d.t.Add(time.Hour * 24 * 7))
}

func (d *Date) PrevWeek() *Date {
	return DateWithTime(d.t.Add(-time.Hour * 24 * 7))
}

func (d *Date) NextMonth() *Date {
	if d.month == 12 {
		return NewDate(d.year+1, 1, d.day)
	}
	return d.Add(0, 1, 0)
}

func (d *Date) PrevMonth() *Date {
	if d.month == 1 {
		return NewDate(d.year-1, 12, d.day)
	}
	return d.Add(0, -1, 0)
}

func (d *Date) Equals(date *Date) bool {
	return d.year == date.year && d.month == date.month && d.day == date.day
}

func (d *Date) Before(date *Date) bool {
	if d.year < date.year {
		return true
	}
	if d.year > date.year {
		return false
	}
	if d.month < date.month {
		return true
	}
	if d.month > date.month {
		return false
	}
	return d.day < date.day
}

func (d *Date) After(date *Date) bool {
	return date.Before(d)
}

func (d *Date) Start() time.Time {
	return time.Date(d.year, time.Month(d.month), d.day, 0, 0, 0, 0, d.t.Location())
}

func (d *Date) End() time.Time {
	return time.Date(d.year, time.Month(d.month), d.day, 23, 59, 59, 999999999, d.t.Location())
}

func (d *Date) IsToday() bool {
	return DateWithTime(time.Now()).Equals(d)
}

func (d *Date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.year, d.month, d.day)
}
