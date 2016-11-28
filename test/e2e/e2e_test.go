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
