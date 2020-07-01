package timex

import (
	"time"
)

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func NumOfYearDays(year int) int {
	if IsLeap(year) {
		return 366
	}
	return 365
}

func GetDayTime(t time.Time) time.Duration {
	return time.Hour*time.Duration(t.Hour()) +
		time.Minute*time.Duration(t.Minute()) +
		time.Second*time.Duration(t.Second()) +
		time.Duration(t.Nanosecond())
}

func IsToday(t time.Time) bool {
	return DateWithTime(t).IsToday()
}

func IsYesterday(t time.Time) bool {
	return DateWithTime(t).IsYesterday()
}

func IsTomorrow(t time.Time) bool {
	return DateWithTime(t).IsTomorrow()
}

func BeginOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local)
}

const (
	Day  = 24 * time.Hour
	Week = 7 * 24 * time.Hour
)

var enWeekdaySymbols = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
var hansWeekdaySymbols = [7]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}

func GetWeekdaySymbol(d int) string {
	d = d % 7
	if IsSimplifiedChinese() {
		return hansWeekdaySymbols[d]
	}
	return enWeekdaySymbols[d]
}
