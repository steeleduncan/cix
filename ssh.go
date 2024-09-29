/*
	ssh.go - Ssh configuration for Cix

	Copyright 2024 Duncan Steele

	Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"fmt"
)

type SshConfiguration struct {
	// The git repository URL
	Remote string
}

var _ RepoSource = &GithubConfiguration{}

func (sc *SshConfiguration) Valid() bool {
	if sc == nil {
		return false
	}

	return sc.Remote != ""
}

func (sc *SshConfiguration) SetStatus(status CiStatus, comment, description, hash string) error {
	// nothing we can do here
	return nil
}

func (sc *SshConfiguration) NixUrl(revision string) string {
	return fmt.Sprintf("git+ssh://%v?rev=%v", sc.Remote, revision)
}

func (sc *SshConfiguration) GitUrl() string {
	return sc.Remote
}
