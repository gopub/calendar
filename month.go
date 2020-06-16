package timex

import "time"

type Month struct {
	Year  int
	Month int
}

func NewMonth(y, m int) *Month {
	if m > 12 {
		m = m % 12
	}
	if m < 1 {
		m = 1
	}
	return &Month{
		Year:  y,
		Month: m,
	}
}

func (m *Month) Begin() time.Time {
	return time.Date(m.Year, time.Month(m.Month), 1, 0, 0, 0, 0, time.Local)
}

func (m *Month) End() time.Time {
	return time.Date(m.Year, time.Month(m.Month+1), 0, 23, 59, 59, 999999999, time.Local)
}

func (m *Month) NumOfDays() int {
	y, mo := m.Year, m.Month
	if mo == 2 {
		if IsLeap(y) {
			return 29
		}
		return 28
	}
	if mo > 7 {
		mo -= 7
	}
	if mo%2 == 0 {
		return 30
	}
	return 31
}

func (m *Month) Equals(month *Month) bool {
	return m.Year == month.Year && m.Month == month.Month
}

func (m *Month) Includes(d *Date) bool {
	return m.Year == d.year && m.Month == d.month
}
