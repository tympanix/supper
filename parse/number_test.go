package parse_test

import (
	"fmt"
	"testing"

	"github.com/Tympanix/supper/parse"
)

func TestPhoneticNumbers(t *testing.T) {
	for i := 1; i <= 50; i++ {
		fmt.Printf("%-4d%-14s\n", i, parse.PhoneticNumber(i))
	}
}
