package timex

import (
	"fmt"
	"time"
)

type Month struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

func NewMonth(y, m int) *Month {
	if m < 0 {
		panic(fmt.Sprintf("timex: month cannot be negative %d", m))
	}
	if m > 12 {
		m = m % 12
	} else if m == 0 {
		m = 12
		y -= 1
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

func (m *Month) Since(month *Month) int {
	return 12*(m.Year-month.Year) + m.Month - month.Month
}

func (m *Month) Add(years, months int) *Month {
	years += m.Year
	months += m.Month
	years += months / 12
	months %= 12
	if months < 0 {
		years -= 1
		months += 12
	}
	return NewMonth(years, months)
}

func CurrentMonth() *Month {
	t := time.Now()
	return NewMonth(t.Year(), int(t.Month()))
}
