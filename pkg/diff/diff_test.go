package diff

import (
	"testing"
	"time"
)

func TestGetLineTime(t *testing.T) {

	line := "Nov 27 15:07:47 hostname1 log file line 5"
	format := "Jan 02 15:04:05"
	date := "Nov 27 15:07:47"

	actual, err := getLineTime(line)
	if err != nil {
		t.Error(err)
	}

	expected, err := time.Parse(format, date)
	if err != nil {
		t.Error(err)
	}
	if actual != expected {
		t.Errorf("Expected time to be %s got %s", expected, actual)
	}
}
