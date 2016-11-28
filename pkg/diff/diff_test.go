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

func TestNewLine(t *testing.T) {
	text := "Nov 27 15:07:47 hostname1 log file line 5"
	line, err := newLine(text)
	if err != nil {
		t.Fatal(err)
	}

	if line.content != text {
		t.Errorf("Line content does not equal text.\nExpected: %s\nGot:%s",
			text, line.content)
	}

	format := "Jan 02 15:04:05"
	date := "Nov 27 15:07:47"
	expected, err := time.Parse(format, date)
	if err != nil {
		t.Error(err)
	}
	if line.time != expected {
		t.Errorf("Line time does not equal text.\nExpected: %v\nGot:%v",
			expected, line.time)
	}
}
