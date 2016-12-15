package e2e

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const logdiff string = "../../_build/logdiff"
const log1 string = "../log1.txt"
const log2 string = "../log2.txt"

// TestOneLog tests reading one log and writing it's output to the console.
// This mainly checks for catastrophic problems, or that someone ran this
// without make build being ran.
func TestOneLog(t *testing.T) {
	args := []string{"--print", log1}
	out, err := exec.Command(logdiff, args...).CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal(err)
	}

	lines := strings.Split(string(out), "\n")

	l, err := os.Open(log1)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	s := bufio.NewScanner(l)

	for i := 0; s.Scan(); i++ {
		actual := s.Text()
		expected := lines[i]
		if expected != actual {
			t.Errorf("%s:\n     Got:\n            %s\nExpected:\n%            s\n",
				log1, expected, actual)
		}
	}
}

func TestTwoLogs(t *testing.T) {
	args := []string{"--print", log1, log2}
	out, err := exec.Command(logdiff, args...).CombinedOutput()
	if err != nil {
		t.Log(string(out))
		t.Fatal(err)
	}

	actual := strings.Split(string(out), "\n")

	expected := []string{"Nov 27 14:33:59 hostname1 log file line 1 ]|[ Nov 27 14:33:59 hostname2 log file line 1",
		"Nov 27 14:34:05 hostname1 log file line 2 ]|[ Nov 27 14:34:05 hostname2 log file line 2",
		"Nov 27 14:34:05 hostname1 log file line 3 ]|[ Nov 27 14:34:05 hostname2 log file line 3",
		" ]|[ Nov 27 15:06:47 hostname2 log file line 4",
		"Nov 27 15:07:47 hostname1 log file line 4 ]|[ Nov 27 15:07:47 hostname2 log file line 5",
		"Nov 27 15:07:47 hostname1 log file line 5 ]|[ Nov 27 15:07:47 hostname2 log file line 6",
		"Nov 27 15:07:47 hostname1 log file line 6 ]|[ ",
		"Nov 27 15:07:47 hostname1 log file line 7 ]|[ ",
		"Nov 27 15:07:47 hostname1 log file line 8 ]|[ ",
		" ]|[ Nov 27 15:07:49 hostname2 log file line 7",
		" ]|[ Nov 27 15:08:49 hostname2 log file line 8"}

	t.Logf("Comparing %s and %s", log1, log2)
	for i, line := range expected {
		if line != actual[i] {
			t.Errorf("\n     Got:\n            '%s'\nExpected:\n            '%s'\n",
				actual[i], line)
		}
	}
}
