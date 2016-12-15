package diff

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func TestScanLine(t *testing.T) {
	f, err := os.Open("../../test/log1.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	lg, err := NewLog(f, f)
	if err != nil {
		t.Error(err)
	}

	i := 1
	for {
		ok := lg.scanLine()
		if !ok {
			break
		}
		if !strings.Contains(lg.currentLine.content, fmt.Sprintf("line %d", i)) {
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

func TestByOldestLines(t *testing.T) {
	tmp1, err := ioutil.TempFile("", "logdiff")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer tmp1.Close()
	defer exec.Command("rm", "-fv", tmp1.Name()).Run()
	tmp2, err := ioutil.TempFile("", "logdiff")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer tmp2.Close()
	defer exec.Command("rm", "-fv", tmp2.Name()).Run()
	l1, err := NewLog(strings.NewReader("Nov 27 14:33:59 hostname1 log file line 1\n"), tmp1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l2, err := NewLog(strings.NewReader("Nov 27 15:07:47 hostname2 log file line 1\n"), tmp2)
	if err != nil {
		t.Fatalf(err.Error())
	}
	logs := []*Log{&l1, &l2}
	err = ByOldestLines(logs...)
	if err != nil {
		t.Fatal(err.Error())
	}

	expected := [][]string{
		{"Nov 27 14:33:59 hostname1 log file line 1", ""},
		{"", "Nov 27 15:07:47 hostname2 log file line 1"},
	}

	for i, l := range logs {
		f, _ := os.Open(l.OutFile.Name())
		s := bufio.NewScanner(f)
		for j := 0; j < len(expected); j++ {
			s.Scan()
			if s.Err() != nil {
				t.Error(s.Err())
			}
			if s.Text() != expected[i][j] {
				t.Errorf("Expected: '%s' \nGot: '%s'", expected[i][j], s.Text())
			}
		}
	}
}
