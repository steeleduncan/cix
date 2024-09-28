/*
   cix.go - Core of the Cix daemon

   Copyright 2024 Duncan Steele

   Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import (
	"fmt"
	"path/filepath"
)

type CiStatus int
const (
	KInProgress CiStatus = iota
	KFailed
	KError
	KSucceeded
)

type Notifier interface {
	SetStatus(status CiStatus, comment, hash string) error
}

type Operation struct {
	Receiver Notifier

	// Repository path
	Repo Repository

	// The hash path
	Hash string
}

func (op Operation) Execute(name string) error {
	if op.Receiver != nil {
		op.Receiver.SetStatus(KInProgress, name, op.Hash)
	}
	ok, err := RunChecks(op.Repo.Path, op.Hash)
	if err != nil {
		if op.Receiver != nil {
			op.Receiver.SetStatus(KError, name, op.Hash)
		}
		return err
	}

	if ok {
		if op.Receiver != nil {
			op.Receiver.SetStatus(KSucceeded, name, op.Hash)
		}
		fmt.Println("Passed!")
	} else {
		if op.Receiver != nil {
			op.Receiver.SetStatus(KFailed, name, op.Hash)
		}
		fmt.Println("Failed!")
	}

	return nil
}

// Perform a single tick
func (c Configuration) GatherNewCommits(varFolder string) ([]Operation, error) {
	if c.Verbose {
		fmt.Println("Gather commits")
	}

	ops := []Operation {}

	for _, repo := range c.Repositories {
		r := Repository {
			Path: filepath.Join(varFolder, "repositories", repo.Identifier()),
		}
		if c.Verbose {
			fmt.Println(" Repository ", r.Path)
		}

		if !r.Exists() {
			if c.Verbose {
				fmt.Println("  Clone ", repo.ResolvedRemote(), repo.Branch)
			}
			if err := r.Clone(repo.ResolvedRemote(), repo.Branch); err != nil {
				return nil, err
			}
		}

		commitsBefore, err := r.ListCommits(repo.Branch)
		if err != nil {
			return nil, err
		}
		if c.Verbose {
			fmt.Println("  ", len(commitsBefore), " commits before fetch")
		}

		if c.Verbose {
			fmt.Println("  Fetch")
		}
		if err := r.Fetch(repo.Branch); err != nil {
			return nil, err
		}

		commitsAfter, err := r.ListCommits(repo.Branch)
		if err != nil {
			return nil, err
		}
		if c.Verbose {
			fmt.Println("  ", len(commitsBefore), " commits after fetch")
		}

		for hash, _ := range commitsAfter {
			_, fnd := commitsBefore[hash]
			if fnd {
				continue
			}
			if c.Verbose {
				fmt.Println("  New commit ", hash)
			}

			op := Operation {
				Repo: r,
				Hash: hash,
				Receiver: nil,
			}
			if repo.Github.Valid() {
				op.Receiver = repo.Github
			}
			ops = append(ops, op)
		}

		op := Operation {
			Repo: r,
			Hash: "e230319a5d9c8b96d445fc863b5084d9b9493c3b",
		}
		if repo.Github.Valid() {
			op.Receiver = repo.Github
		}
		ops = append(ops, op)
	}

	return ops, nil
}

func (c Configuration) Tick() error {
	varFolder := filepath.Join(c.Var, "v1")

	ops, err := c.GatherNewCommits(varFolder)
	if err != nil {
		return err
	}	

	for _, op := range ops {
		err := op.Execute(c.ResolvedName())
		if err != nil {
			return err
		}
	}

	return nil
}
