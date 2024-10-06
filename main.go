/*
main.go - Cli interface for go

# Copyright 2024 Duncan Steele

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/steeleduncan/cix/version"
)

func usage() error {
	fmt.Println(`cix ` + version.Version() + ` <config.json>`)
	return nil
}

func errMain() error {
	if len(os.Args) != 2 {
		return usage()
	}

	blob, err := os.ReadFile(os.Args[1])
	if err != nil {
		return fmt.Errorf("Failed to read config at %v", os.Args[1])
	}

	c := Configuration{}
	err = json.Unmarshal(blob, &c)
	if err != nil {
		return fmt.Errorf("Bad json: %v", err)
	}
	c.Var = os.ExpandEnv(c.Var)

	for {
		err := c.Tick()
		if err != nil {
			fmt.Println("error: ", err)
		}

		time.Sleep(3 * time.Minute)
	}
}

func main() {
	err := errMain()
	if err != nil {
		fmt.Println("fatal: ", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
