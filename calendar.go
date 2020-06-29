package timex

type Repeat int

const (
	Never Repeat = iota
	Daily
	Weekly
	Monthly
	Yearly
)
