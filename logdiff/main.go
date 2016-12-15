/*
logdiff compares log files based on time
Copyright © 2016 Hayley Swimelar

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the "Software"),
to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense,
and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included
in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hayswim/logdiff/pkg/diff"
	"github.com/hayswim/logdiff/pkg/display"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func run() error {
	print := flag.Bool("print", false, "Print Diff To stdout")
	flag.Parse()

	if !*print {
		panic("Not Implemented")
	}

	args := flag.Args()
	if len(args) < 1 {
		return errors.New("Expected a file")
	}

	var logs []*diff.Log
	for _, file := range args {
		file, err := os.Open(file)
		if err != nil {
			return err
		}
		defer file.Close()

		tmp, err := ioutil.TempFile("", "logdiff")
		if err != nil {
			return err
		}
		defer tmp.Close()

		l, err := diff.NewLog(file, tmp)
		if err != nil {
			return err
		}
		logs = append(logs, &l)
	}
	err := diff.ByOldestLines(logs...)
	if err != nil {
		return err
	}

	display.Print(" ]|[ ", logs...)
	return nil
}
