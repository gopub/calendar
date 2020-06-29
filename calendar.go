package timex

import (
	"fmt"
	"os"
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
var zhHansRepeats = []string{"不重复", "每日", "每年", "每月", "每年"}

func (r Repeat) String() string {
	if r < Never || r > Yearly {
		return fmt.Sprint(int(r))
	}
	lang := getLang()
	if isHans(lang) {
		return zhHansRepeats[r]
	}
	return enRepeats[r]
}

func getLang() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en_US"
	}
	return lang
}

var hansLangList = []string{"zh_CN", "zh_SG"}

func isHans(lang string) bool {
	for _, s := range hansLangList {
		if strings.Contains(lang, s) {
			return true
		}
	}
	return false
}
