package parse

import (
	"fmt"
	"strings"
)

var singles = []string{
	"First",
	"Second",
	"Third",
	"Fourth",
	"Fifth",
	"Sixth",
	"Seventh",
	"Eighth",
	"Ninth",
	"Tenth",
	"Eleventh",
	"Twelfth",
	"Thirteenth",
	"Fourteenth",
	"Fifteenth",
	"Sixteenth",
	"Seventeenth",
	"Eighteenth",
	"Nineteenth",
}

var tens = []string{
	"Twenty",
	"Thirty",
	"Fourty",
	"Fifty",
	"Sixty",
	"Seventy",
	"Eighty",
	"Ninety",
}

// PhoneticNumber converts an integer into a human readable string
func PhoneticNumber(num int) string {
	if num <= 0 {
		return ""
	} else if num < 20 {
		return singles[num-1]
	} else if num%10 == 0 {
		base := strings.TrimSuffix(tens[num/10-2], "ty")
		return base + "tieth"
	} else if num < 100 {
		return fmt.Sprintf("%s-%s", tens[num/10-2], singles[num%10-1])
	} else {
		return ""
	}
}
