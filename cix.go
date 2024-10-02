/*
cix.go - Core of the Cix daemon

# Copyright 2024 Duncan Steele

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

type RepoSource interface {
	// True if this is useable (safe against nils)
	Valid() bool

	// Set the commit status (using the forge's commitstatus api)
	SetStatus(status CiStatus, comment, description, hash string) error

	// Nix url, used for printing in the status description
	NixUrl(revision string) string

	// Git url for cloning
	GitUrl() string
}

type Operation struct {
	Source RepoSource

	// Repository path
	Repo Repository

	// The hash path
	Hash string
}

func (c Configuration) Execute(op Operation) error {
	name := c.ResolvedName()

	description := GetDescription(op.Hash, op.Source)

	if op.Source != nil {
		op.Source.SetStatus(KInProgress, name, "", op.Hash)
	}
	ok, err := c.RunChecks(op.Repo.Path, op.Hash)
	if err != nil {
		if op.Source != nil {
			op.Source.SetStatus(KError, name, description, op.Hash)
		}
		return err
	}

	if ok {
		if op.Source != nil {
			op.Source.SetStatus(KSucceeded, name, description, op.Hash)
		}
		fmt.Println("Passed!")
	} else {
		if op.Source != nil {
			op.Source.SetStatus(KFailed, name, description, op.Hash)
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

	ops := []Operation{}

	for _, repo := range c.Repositories {
		r := Repository{
			Path: filepath.Join(varFolder, repo.Identifier()),
		}
		if c.Verbose {
			fmt.Println(" Repository ", r.Path)
		}

		source := repo.Source()

		if !r.Exists() {
			if c.Verbose {
				fmt.Println("  Clone ", source.GitUrl(), repo.Branch)
			}
			if err := r.Clone(source.GitUrl(), repo.Branch); err != nil {
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

			op := Operation{
				Repo:   r,
				Hash:   hash,
				Source: repo.Source(),
			}
			ops = append(ops, op)
		}
	}

	return ops, nil
}

func (c Configuration) Validate() error {
	for i, repo := range c.Repositories {
		if repo.Source() == nil {
			return fmt.Errorf("Invalid configuration: missing repository source (%v)", i)
		}
	}
	return nil
}

func (c Configuration) Tick() error {
	if err := c.Validate(); err != nil {
		return err
	}

	varFolder := filepath.Join(c.Var, "v1")

	ops, err := c.GatherNewCommits(varFolder)
	if err != nil {
		return err
	}

	for _, op := range ops {
		err := c.Execute(op)
		if err != nil {
			return err
		}
	}

	return nil
}
