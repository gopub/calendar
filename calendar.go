package calendar

import "time"

const (
	EndOfDay = time.Hour*24 - time.Nanosecond
)

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func LenOfYear(year int) int {
	if IsLeap(year) {
		return 366
	}
	return 365
}

func LenOfMonth(year int, month int) int {
	if month == 2 {
		if IsLeap(year) {
			return 29
		}
		return 28
	}
	if month > 7 {
		month -= 7
	}
	if month%2 == 0 {
		return 30
	}
	return 31
}

func BeginOfMonth(year int, month int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
}

func EndOfMonth(year, month int) time.Time {
	return time.Date(year, time.Month(month), LenOfMonth(year, month), 23, 59, 59, 999999999, time.Local)
}

func TimeInDay(t time.Time) time.Duration {
	return time.Hour*time.Duration(t.Hour()) +
		time.Minute*time.Duration(t.Minute()) +
		time.Second*time.Duration(t.Second()) +
		time.Duration(t.Nanosecond())
}

func IsToday(t time.Time) bool {
	return DateWithTime(t).IsToday()
}
