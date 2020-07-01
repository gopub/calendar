package timex

import (
	"fmt"
	"time"
)

// Time is for mobile
type Time struct {
	t time.Time
}

func NewTime() *Time {
	return &Time{t: time.Now()}
}

func TimeWithUnix(sec int64) *Time {
	return &Time{t: time.Unix(sec, 0)}
}

func (t *Time) Unix() int64 {
	return t.t.Unix()
}

func (t *Time) Year() int {
	return t.t.Year()
}

func (t *Time) Month() int {
	return int(t.t.Month())
}

func (t *Time) Day() int {
	return t.t.Day()
}

func (t *Time) Hour() int {
	return t.t.Hour()
}

func (t *Time) Minute() int {
	return t.t.Minute()
}

func (t *Time) Weekday() int {
	return int(t.t.Weekday())
}

func (t *Time) Date() *Date {
	return NewDate(t.Year(), t.Month(), t.Day())
}

func (t *Time) DayMinutes() int {
	rt := t.t.Add(time.Nanosecond).Round(time.Minute)
	return rt.Hour()*60 + rt.Minute()
}

func (t *Time) Before(tm *Time) bool {
	return t.t.Before(tm.t)
}

func (t *Time) After(tm *Time) bool {
	return t.t.After(tm.t)
}

func (t *Time) Equals(tm *Time) bool {
	return t.t.Equal(tm.t)
}

func (t *Time) AddDate(years, months, days int) *Time {
	return &Time{t: t.t.AddDate(years, months, days)}
}

func (t *Time) AddTime(hours, minutes int) *Time {
	return &Time{
		t: t.t.Add(time.Hour*time.Duration(hours) + time.Minute*time.Duration(minutes)),
	}
}

func (t *Time) AddHours(hours int) *Time {
	return &Time{
		t: t.t.Add(time.Hour * time.Duration(hours)),
	}
}

func (t *Time) AddMinutes(minutes int) *Time {
	return &Time{
		t: t.t.Add(time.Minute * time.Duration(minutes)),
	}
}

func (t *Time) AddSeconds(sec int) *Time {
	return &Time{
		t: t.t.Add(time.Second * time.Duration(sec)),
	}
}

func (t *Time) AddNanos(nano int) *Time {
	return &Time{
		t: t.t.Add(time.Nanosecond * time.Duration(nano)),
	}
}

func (t *Time) IsBeginOfDay() bool {
	return GetDayTime(t.t) == 0
}

func (t *Time) IsEndOfDay() bool {
	return GetDayTime(t.t) == time.Hour*24-time.Nanosecond
}

func (t *Time) RoundHour() *Time {
	return &Time{
		t: t.t.Round(time.Hour),
	}
}

func (t *Time) TimeText() string {
	return t.timeText("%d:%02d")
}

func (t *Time) TimeTextWithZero() string {
	return t.timeText("%02d:%02d")
}

func (t *Time) TimeTextWithSpace() string {
	return t.timeText("%2d:%02d")
}

func (t *Time) timeText(format string) string {
	if IsSimplifiedChinese() {
		if t.IsEndOfDay() {
			return "晚上12:00"
		}
		h, m := t.Hour(), t.Minute()
		switch {
		case h == 0:
			return fmt.Sprintf("00:%02d", m)
		case h < 12:
			return fmt.Sprintf("上午"+format, h, m)
		case h == 12:
			return fmt.Sprintf("中午"+format, h, m)
		default:
			return fmt.Sprintf("下午"+format, h-12, m)
		}
	}
	if t.IsEndOfDay() {
		return "Midnight"
	}
	h, m := t.Hour(), t.Minute()
	switch {
	case h == 0:
		return fmt.Sprintf(format+"AM", 12, m)
	case h < 12:
		return fmt.Sprintf(format+"AM", h, m)
	case h == 12:
		return fmt.Sprintf(format+"PM", h, m)
	default:
		return fmt.Sprintf(format+"PM", h-12, m)
	}
}

func (t *Time) RelativeDateTimeText() string {
	return fmt.Sprintf("%s %s", t.Date().ShortRelativeText(), t.TimeTextWithZero())
}

func NewRangeT(begin, end *Time) *Range {
	return NewRange(begin.t, end.t)
}

func NewRangeInDay(d *Date) *Range {
	return NewRange(d.Begin(), d.End().Add(time.Nanosecond))
}

func (r *Range) SetT(begin, end *Time) {
	r.Set(begin.t, end.t)
}

func (r *Range) SetBeginT(t *Time) {
	r.SetBegin(t.t)
}

func (r *Range) SetEndT(t *Time) {
	r.SetEnd(t.t)
}

func (r *Range) BeginT() *Time {
	return &Time{t: r.begin}
}

func (r *Range) EndT() *Time {
	return &Time{t: r.end}
}

func (d *Date) BeginT() *Time {
	return &Time{t: d.Begin()}
}

func (d *Date) EndT() *Time {
	return &Time{t: d.End()}
}
