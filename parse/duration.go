package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var timeDuration = map[string]time.Duration{
	"s": time.Second,
	"m": time.Minute,
	"h": time.Hour,
	"d": time.Hour * 24,
}

var letterRegexp = regexp.MustCompile(`[a-zA-Z]+`)
var digitRegexp = regexp.MustCompile(`[0-9]+`)

func removeEmptyString(xls []string) []string {
	ls := make([]string, 0)
	for _, s := range xls {
		if len(s) > 0 {
			ls = append(ls, s)
		}
	}
	return ls
}

// Duration takes a textual representation of a duration (e.g. 2h38m) and
// returns the parsed duration from the time package. An error is returned if
// the input was invalid (i.e. could not be parsed into a duration)
func Duration(str string) (time.Duration, error) {
	t := time.Duration(0)
	vals := removeEmptyString(letterRegexp.Split(str, -1))
	mods := removeEmptyString(digitRegexp.Split(str, -1))

	if len(vals) != len(mods) {
		return 0, fmt.Errorf("could not parse time format")
	}

	for i := range vals {
		num, err := strconv.Atoi(vals[i])
		if err != nil {
			return 0, fmt.Errorf("could not parse time format")
		}
		mod, exists := timeDuration[mods[i]]
		if !exists {
			return 0, fmt.Errorf("unknown time format specifier: %v", mods[i])
		}
		t += time.Duration(num) * mod
	}
	return t, nil
}
