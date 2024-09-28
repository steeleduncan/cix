/*
   git.go - Git tools for cix

   Copyright 2024 Duncan Steele

   Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import (
	"os"
	"io"
	"fmt"
	"errors"
	"strings"
	"os/exec"
	"path/filepath"
)

// Check if a commit is (likely) valid
// This is for filtering out shell junk
func VerifyCommit(s string) bool {
	if len(s) != 40 {
		return false
	}

	return true
}

// Main repo object
// TODO add branch as we only keep a single branch per repo
type Repository struct {
	Path string
}

// Check for existence
func (r Repository) Exists() bool {
	_, err := os.Stat(filepath.Join(r.Path, "config"))
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return true
	}
}

// Checkout the repo to the given path
func (r Repository) CheckoutTo(path, branch string) error {
	cmd := exec.Command("git", "checkout", branch, ".")
	cmd.Dir = r.Path
	cmd.Env = []string { fmt.Sprintf("GIT_WORK_TREE=%v", path) }
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Checkout failed for %v / %v", r.Path, branch)
	}

	return nil
}

// Return a set of all commits in the repository
func (r Repository) ListCommits(branch string) (map[string]bool, error) {
	cmd := exec.Command("git", "rev-list", branch)
	cmd.Dir = r.Path

	so, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("Failed to create stdout pipe")
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("Failed to run git command to list commits")
	}

	slurp, _ := io.ReadAll(so)
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("Failed running git command to list commits")
	}

	ret := map[string]bool {}
	hashes := strings.Split(string(slurp), "\n")
	for _, hash := range hashes {
		if hash == "" {
			continue
		}

		if !VerifyCommit(hash) {
			return nil, fmt.Errorf("Did not understand hash: '%v'", hash)
		}
		ret[hash] = true
	}
	return ret, nil
}

// Fetch all new commits
func (r Repository) Fetch(branch string) error {
	cmd := exec.Command("git", "fetch", "origin", branch + ":" + branch)
	cmd.Dir = r.Path

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Fetch failed for %v / %v", r.Path, branch)
	}

	return nil
}

// Clone a new repo
func (r Repository) Clone(remote, branch string) error {
	// TODO ideally we'd clone to a temporary working path
	name := filepath.Base(r.Path)
	parent := filepath.Dir(r.Path)
	os.MkdirAll(parent, 0777)
	
	cmd := exec.Command("git", "clone", "--bare", "-b", branch, remote, name)
	cmd.Dir = parent
	so, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to create stdout pipe")
	}
	if err := cmd.Start(); err != nil {
		os.RemoveAll(r.Path)
		return fmt.Errorf("Clone failed for %v / %v", remote, branch)
	}

	serr, _ := io.ReadAll(so)
	if err := cmd.Wait(); err != nil {
		fmt.Println(string(serr))
		os.RemoveAll(r.Path)
		return fmt.Errorf("Clone failed for %v / %v", remote, branch)
	}

	return nil
}
