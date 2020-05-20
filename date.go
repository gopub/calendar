package calendar

import "time"

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

func DateWithUnixSeconds(seconds int64) *Date {
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
	return NewDate(d.year+years, d.month+months, d.day+days)
}

func (d *Date) Next() *Date {
	return d.Add(0, 0, 1)
}

func (d *Date) Prev() *Date {
	return d.Add(0, 0, -1)
}

func (d *Date) Equals(date *Date) bool {
	return d.year == date.year && d.month == date.month && d.day == date.day
}

func (d *Date) Before(date *Date) bool {
	return d.year < date.year || d.month < date.month || d.day < date.day
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
