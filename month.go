package timex

import (
	"fmt"
	"sync"
	"time"
)

var monthWeeksNum = &sync.Map{}

type Month struct {
	Year  int `json:"year"`
	Month int `json:"month"`

	// [4,7]*7 matrix
	mu       sync.Mutex
	calendar [][7]*Date
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

func (m *Month) NumOfWeeks() int {
	if m.calendar == nil {
		_ = m.GetCalendarDate(1, 1)
	}
	return len(m.calendar)
}

func (m *Month) Before(mo *Month) bool {
	if m.Year == mo.Year {
		return m.Month < mo.Month
	}
	return m.Year < mo.Year
}

func (m *Month) After(mo *Month) bool {
	return mo.Before(m)
}

// GetCalendarDate: week is [1, NumOfWeeks], day is [1, 7]
func (m *Month) GetCalendarDate(week, day int) *Date {
	week -= 1
	day -= 1
	if m.calendar == nil {
		m.mu.Lock()
		if m.calendar == nil {
			first := int(m.Begin().Weekday())
			numOfDays := m.NumOfDays()
			days := numOfDays - (7 - first)
			lines := 1 + days/7
			if days%7 != 0 {
				lines += 1
			}
			m.calendar = make([][7]*Date, lines)
			last := first + numOfDays - 1
			for i := 0; i < lines; i++ {
				for j := 0; j < 7; j++ {
					offset := i*7 + j
					if offset < first || offset > last {
						m.calendar[i][j] = nil
					} else {
						m.calendar[i][j] = NewDate(m.Year, m.Month, offset-first+1)
					}
				}
			}
		}
		m.mu.Unlock()
	}
	if week < 0 || week >= len(m.calendar) {
		return nil
	}
	if day < 0 || day >= 7 {
		return nil
	}
	return m.calendar[week][day]
}

func (m *Month) Date(day int) *Date {
	return NewDate(m.Year, m.Month, day)
}

func (m *Month) String() string {
	return fmt.Sprintf("%d-%d", m.Year, m.Month)
}

func CurrentMonth() *Month {
	t := time.Now()
	return NewMonth(t.Year(), int(t.Month()))
}
