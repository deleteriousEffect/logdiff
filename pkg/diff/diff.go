package diff

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

import "log"

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

func getLineTime(s string) (time.Time, error) {
	format, err := findTimeFormat(s)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, nil
	}
	lineTime, err := time.Parse(format, s[:len(format)])
	if err != nil {
		fmt.Println(err)
	}
	return lineTime, nil
}

func findTimeFormat(s string) (string, error) {
	//TODO: Make this for real.
	return "Jan 02 15:04:05", nil
}
