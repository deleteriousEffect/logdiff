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
		l.time = time.Time{}
		return nil
	}
	l.time = lineTime
	return nil
}

// Log represents a log.
type Log struct {
	inFile      io.Reader
	OutFile     *os.File
	currentLine line
}

func NewLog(r io.Reader, f *os.File) (Log, error) {
	return Log{r, f, line{}}, nil
}

func (lg *Log) scanLine() bool {
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

func (lg *Log) writeDiff(t time.Time) (int, error) {
	if lg.currentLine.time.Equal(t) || lg.currentLine.time.Before(t) {
		n, err := fmt.Fprintf(lg.OutFile, lg.currentLine.content)
		lg.currentLine = line{}
		if err != nil {
			return n, err
		}
	}
	n, err := fmt.Fprintln(lg.OutFile, "")
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
func ByOldestLines(logs ...*Log) error {
	_ = "breakpoint"
	defer func() {
		for _, lg := range logs {
			lg.OutFile.Sync()
			lg.OutFile.Close()
		}
	}()
	scanners := make([]*bufio.Scanner, len(logs))
	for i, lg := range logs {
		scanners[i] = bufio.NewScanner(lg.inFile)
	}
	for {
		oldestTime := time.Now()
		seenEnd := 0
		for i, lg := range logs {
			if lg.currentLine.content == "" {
				ok := scanners[i].Scan()
				if !ok {
					seenEnd++
					if seenEnd >= len(logs) {
						return nil
					}
					continue
				}
				t := scanners[i].Text()
				ln, err := newLine(t)
				if err != nil {
					return err
				}
				lg.currentLine = ln
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
