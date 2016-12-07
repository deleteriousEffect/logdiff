package diff

import (
	"fmt"
	"os"
	"reflect"
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
	l, err := newLine(text)
	if err != nil {
		t.Fatal(err)
	}

	if l.content != text {
		t.Errorf("Line content does not equal text.\nExpected: %s\nGot:%s",
			text, l.content)
	}

	format := "Jan 02 15:04:05"
	date := "Nov 27 15:07:47"
	expected, err := time.Parse(format, date)
	if err != nil {
		t.Error(err)
	}
	if l.time != expected {
		t.Errorf("Line time does not equal text.\nExpected: %v\nGot:%v",
			expected, l.time)
	}
}

func TestLineReader(t *testing.T) {
	f, err := os.Open("../../test/log1.txt")
	if err != nil {
		t.Fatal(err)
	}
	lg, err := newLog(f)
	defer lg.file.Close()
	if err != nil {
		t.Error(err)
	}

	i := 1
	for {
		l, ok := lg.popLine()
		if !ok {
			break
		}
		if !strings.Contains(l.content, fmt.Sprintf("line %d", i)) {
			t.Errorf("Expected to read line %d", i)
		}
		i++
	}
}

func TestOldestTime(t *testing.T) {
	times := []line{line{"", time.Now()}, line{"", time.Now()}, line{"", time.Time{}}}

	actual := oldestTime(times...)
	expected := time.Time{}

	if expected != actual {
		t.Errorf("Expected: %s, Got: %s", expected, actual)
	}
}

func TestOldestLines(t *testing.T) {
	var linesTests = []struct {
		in  []line
		out []string
	}{
		{[]line{line{"line1", time.Now()}}, []string{"line1"}},
		{[]line{line{"line1", time.Now()}, line{"line2", time.Time{}}}, []string{"", "line2"}},
	}
	for _, lt := range linesTests {
		lines := oldestLines(lt.in...)
		if !reflect.DeepEqual(lines, lt.out) {
			t.Errorf("Expected: %s, Got: %s", lt.out, lines)
		}
	}
}

//func TestByOldestLines(t *testing.T) {
//	var oldestLinesTests = []struct {
//		in  []*bufio.Scanner
//		out [][]string
//	}{
//		{[]*bufio.Scanner{
//			bufio.NewScanner(strings.NewReader("Nov 27 14:33:59 hostname1 log file line 1\n")),
//			bufio.NewScanner(strings.NewReader("Nov 27 15:07:47 hostname2 log file line 1\n"))},
//			[][]string{
//				{"Nov 27 14:33:59 hostname1 log file line 1\n", "\n"},
//				{"\n", "Nov 27 15:07:47 hostname2 log file line 1\n"}}},
//	}
//	for _, blt := range oldestLinesTests {
//		files, err := ByOldestLines(blt.in)
//		if err != nil {
//			t.Error(err)
//		}
//		if true {
//			t.Error(files)
//		}
//	}
//}
