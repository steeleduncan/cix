/*
config.go - Configuration of Cix

# Copyright 2024 Duncan Steele

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"crypto/sha256"
	"fmt"
)

type RepositoryConfiguration struct {
	// (optional) BB config block
	Bitbucket *BitbucketConfiguration

	// (optional) GH config block
	Github *GithubConfiguration

	// (optional) Ssh config block
	Ssh *SshConfiguration

	// The branch to check
	Branch string
}

func (rc RepositoryConfiguration) Source() RepoSource {
	if rc.Bitbucket.Valid() {
		return rc.Bitbucket
	}

	if rc.Github.Valid() {
		return rc.Github
	}

	if rc.Ssh.Valid() {
		return rc.Ssh
	}

	return nil
}

func (rc RepositoryConfiguration) Identifier() string {
	h := sha256.New()
	h.Write([]byte(rc.Source().GitUrl()))
	h.Write([]byte(rc.Branch))
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Configuration struct {
	// Path to our data folder
	Var string

	// A name for this runner
	Name string

	// When true we print a lot
	Verbose bool

	// Timeout in seconds for a build
	Timeout int

	// Path to nix
	NixPath string

	// various git repos
	Repositories []RepositoryConfiguration
}


func (rc Configuration) ResolvedTimeout() int {
	if rc.Timeout == 0 {
		return 15 * 60
	}

	return rc.Timeout
}

func (rc Configuration) ResolvedNixPath() string {
	if rc.NixPath == "" {
		// hope it is in the path
		return "nix"
	}

	return rc.NixPath
}

func (rc Configuration) ResolvedName() string {
	if rc.Name == "" {
		return "Cix"
	}

	return fmt.Sprintf("%v (Cix)", rc.Name)
}
