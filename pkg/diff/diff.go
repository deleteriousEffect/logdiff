package diff

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type line struct {
	// The txt of the line, including timestamp.
	content string
	// The parsed time of the line.
	time time.Time
}

func newLine(s string) (line, error) {
	l := line{s, time.Time{}}
	err := l.setTime()
	if err != nil {
		return line{}, err
	}
	return l, nil
}

func (l *line) setTime() error {
	format, err := findTimeFormat(l.content)
	if err != nil {
		return err
	}
	lineTime, err := time.Parse(format, l.content[:len(format)])
	if err != nil {
		return err
	}
	l.time = lineTime
	return nil
}

type log struct {
	inFile      io.Reader
	outFile     *os.File
	currentLine line
}

func newLog(r io.Reader, f *os.File) (log, error) {
	return log{r, f, line{}}, nil
}

func (lg *log) scanLine() bool {
	s := bufio.NewScanner(lg.inFile)
	ok := s.Scan()
	if !ok {
		return ok
	}

	t := s.Text()
	ln, err := newLine(t)
	if err != nil {
		return false
	}
	lg.currentLine = ln
	return ok
}

func (lg *log) writeDiff(t time.Time) (int, error) {
	if lg.currentLine.time.Equal(t) || lg.currentLine.time.Before(t) {
		n, err := fmt.Fprintf(lg.outFile, lg.currentLine.content)
		lg.currentLine = line{}
		if err != nil {
			return n, err
		}
	}
	n, err := fmt.Fprintln(lg.outFile, "")
	if err != nil {
		return n, err
	}
	return n, nil
}

func findTimeFormat(s string) (string, error) {
	//TODO: Make this for real.
	return "Jan 02 15:04:05", nil
}

func oldestTime(lines ...line) time.Time {
	oldest := time.Now()

	for _, l := range lines {
		if l.time.Before(oldest) {
			oldest = l.time
		}
	}
	return oldest
}

func oldestLines(lines ...line) []string {
	oldest := []string{}
	ot := oldestTime(lines...)
	for _, l := range lines {
		if l.time.Equal(ot) || l.time.Before(ot) {
			oldest = append(oldest, l.content)
		} else {
			oldest = append(oldest, "")
		}
	}
	return oldest
}

// ByOldestLines diffs logs based on the time each line was logged and writes
// the results to the logs outFile.
func ByOldestLines(logs ...*log) error {
	defer func() {
		for _, lg := range logs {
			lg.outFile.Sync()
			lg.outFile.Close()
		}
	}()
	for {
		oldestTime := time.Now()
		seenEnd := 0
		for _, lg := range logs {
			if lg.currentLine.content == "" {
				ok := lg.scanLine()
				if !ok {
					seenEnd++
				}
			}
			if seenEnd >= len(logs) {
				return nil
			}
			if lg.currentLine.time.Equal(oldestTime) || lg.currentLine.time.Before(oldestTime) {
				oldestTime = lg.currentLine.time
			}
		}
		for _, lg := range logs {
			n, err := lg.writeDiff(oldestTime)
			if err != nil {
				return err
			}
			if n == 0 {
				return errors.New("Nothing written to outFile")
			}
		}
	}
}
