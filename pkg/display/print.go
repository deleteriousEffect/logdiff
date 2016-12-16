package display

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hayswim/logdiff/pkg/diff"
)

// Print writes to standard out.
func Print(sep string, logs ...*diff.Log) {
	defer func() {
		for _, lg := range logs {
			lg.OutFile.Close()
		}
	}()
	var scanners []*bufio.Scanner
	for _, lg := range logs {
		f, _ := os.Open(lg.OutFile.Name())
		scanners = append(scanners, bufio.NewScanner(f))
	}
	for {
		seenEnd := 0
		currentLines := make([]string, len(logs))
		for i, s := range scanners {
			ok := s.Scan()
			if !ok {
				seenEnd++
			}
			if seenEnd >= len(logs) {
				return
			}
			currentLines[i] = s.Text()
		}
		fmt.Println(strings.Join(currentLines, sep))
	}
}
