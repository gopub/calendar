package timex

import (
	"fmt"
	"strings"
)

type Repeat int

const (
	Never Repeat = iota
	Daily
	Weekly
	Monthly
	Yearly
)

var enRepeats = []string{"Never", "Daily", "Weekly", "Monthly", "Yearly"}
var zhHansRepeats = []string{"不重复", "每天", "每年", "每月", "每年"}

func (r Repeat) IsValid() bool {
	switch r {
	case Never, Daily, Weekly, Monthly, Yearly:
		return true
	default:
		return false
	}
}

func (r Repeat) String() string {
	if r < Never || r > Yearly {
		return fmt.Sprint(int(r))
	}
	if IsSimplifiedChinese() {
		return zhHansRepeats[r]
	}
	return enRepeats[r]
}

type langT int

const (
	english = iota
	simplifiedChinese
	traditionalChinese
)

var lang langT = english

func SetLang(l string) {
	l = strings.ToLower(l)
	strings.Replace(l, "-", "_", -1)
	switch {
	case strings.Contains(l, "hans"):
		lang = simplifiedChinese
	case strings.Contains(l, "zh_cn"):
		lang = simplifiedChinese
	case strings.Contains(l, "zh_sg"):
		lang = simplifiedChinese
	}
}

func IsSimplifiedChinese() bool {
	return lang == simplifiedChinese
}

func IsTraditionalChinese() bool {
	return lang == traditionalChinese
}
