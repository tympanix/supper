package parse

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	d, err := Duration("2h3m")

	if err != nil {
		t.Fatal(err)
	}

	res := 2*time.Hour + 3*time.Minute
	if d != res {
		t.Fatal("unexpected duration value")
	}
}
