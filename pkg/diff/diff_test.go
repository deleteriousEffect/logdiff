package diff

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestSetTime(t *testing.T) {

	l := line{"Nov 27 15:07:47 hostname1 log file line 5", time.Time{}}
	format := "Jan 02 15:04:05"
	date := "Nov 27 15:07:47"

	err := l.setTime()
	if err != nil {
		t.Error(err)
	}

	expected, err := time.Parse(format, date)
	if err != nil {
		t.Error(err)
	}
	if l.time != expected {
		t.Errorf("Expected time to be %s got %s", expected, l.time)
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

func TestLineReader(t *testing.T) {
	f, err := os.Open("../../test/log1.txt")
	if err != nil {
		t.Fatal(err)
	}
	nextLine := lineReader(f)

	i := 1
	for {
		l, ok := nextLine()
		if !ok {
			break
		}
		if !strings.Contains(l.content, fmt.Sprintf("line %d", i)) {
			t.Errorf("Expected to read line %d", i)
		}
		i++
	}
}
