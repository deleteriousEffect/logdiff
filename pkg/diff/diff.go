package diff

import (
	"bufio"
	"io/ioutil"
	"os"
	"time"
)

import "log"

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

// ByOldestLines diffs files based on the time each line was logged and returns
// a []string of temporary file names where the returns were written.
func ByOldestLines(f ...*os.File) ([]*os.File, error) {

	numFiles := len(f)

	// Initialize tempfiles, scanners, and buffers.
	tempFiles := make([]*os.File, numFiles)
	scanners := make([]*bufio.Scanner, numFiles)
	buffers := make([]string, numFiles)
	for i, file := range f {
		temp, err := ioutil.TempFile("/tmp", file.Name())
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		tempFiles[i] = temp
		scanners[i] = bufio.NewScanner(file)
	}

	reachedEnd := 0
	for {
		for i := 0; i < numFiles; i++ {
			if buffers[i] == "" {
				ok := scanners[i].Scan()
				if !ok {
					buffers[i] = "\n"
					reachedEnd++
					continue
				}
				buffers[i] = scanners[i].Text()
			}
		}

		if reachedEnd >= numFiles {
			break
		}
	}
	return tempFiles, nil
}
