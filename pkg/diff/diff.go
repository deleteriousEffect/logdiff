package diff

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func findTimeFormat(s string) (string, error) {
	//TODO: Make this for real.
	return "Jan 02 15:04:05", nil
}

func lineReader(r io.Reader) func() (line, bool) {
	scanner := bufio.NewScanner(r)
	return func() (line, bool) {
		ok := scanner.Scan()
		if !ok {
			return line{}, ok
		}

		t := scanner.Text()
		l, err := newLine(t)
		if err != nil {
			log.Fatal(err)
		}
		return l, ok
	}
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

// ByOldestLines diffs files based on the time each line was logged and returns
// a []string of temporary file names where the returns were written.
func ByOldestLines(f ...io.ReadWriter) ([]*os.File, error) {

	numFiles := len(f)

	// Initialize tempfiles, scanners, and lines.
	tempFiles := make([]*os.File, numFiles)
	lineReaders := make([]func() (line, bool), numFiles)
	lines := make([]line, numFiles)
	for i, file := range f {
		temp, err := ioutil.TempFile("/tmp", fmt.Sprintf("logdiff_tmp%d", i))
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		tempFiles[i] = temp
		lineReaders[i] = lineReader(file)
	}

	for {
		reachedEnd := 0
		// Create a new line from each file if the line content is blank.
		for i, l := range lines {
			if l.content == "" {
				l, ok := lineReaders[i]()
				// If we can't read it, assume we've reached the end of the file.
				if !ok {
					reachedEnd++
				}
				lines[i] = l
			}
		}

		// If a file starts with the oldest timestamp, write it to the tempfiles.
		// Otherwise, write a newline.
		for i, f := range tempFiles {
			oldest := oldestTime(lines...)
			if lines[i].time.Equal(oldest) {
				_, err := f.WriteString(lines[i].content)
				if err != nil {
					return nil, err
				}
				lines[i] = line{}
				continue
			}
			_, err := f.WriteString("\n")
			if err != nil {
				return nil, err
			}

		}
		// If we've reached the end of every file, we're done.
		if reachedEnd >= numFiles {
			return tempFiles, nil
		}
	}
}
