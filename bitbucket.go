/*
bitbucket.go - Bitbucket tools for Cix

# Copyright 2024 Duncan Steele

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type BitbucketConfiguration struct {
	Workspace  string
	Repository string
	Token      string
}

var _ RepoSource = &BitbucketConfiguration{}

func (bc *BitbucketConfiguration) NixUrl(revision string) string {
	return fmt.Sprintf("git+ssh://git@bitbucket.org:%v/%v?rev=%v", bc.Workspace, bc.Repository, revision)
}

func (bc *BitbucketConfiguration) GitUrl() string {
	return fmt.Sprintf("git@bitbucket.org:%v/%v", bc.Workspace, bc.Repository)
}

func (bc *BitbucketConfiguration) Valid() bool {
	if bc == nil {
		return false
	}

	return bc.Workspace != "" && bc.Repository != ""
}

func (bc *BitbucketConfiguration) SetStatus(status CiStatus, comment, description, hash string) error {
	if bc.Token == "" {
		return nil
	}

	method := "PUT"

	st := "error"
	optionalKey := comment
	switch status {
	case KError:
		st = "STOPPED"

	case KInProgress:
		method = "POST"
		st = "INPROGRESS"
		optionalKey = ""

	case KFailed:
		st = "FAILED"

	case KSucceeded:
		st = "SUCCESSFUL"
	}

	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%v/%v/commit/%v/statuses/build/%v", bc.Workspace, bc.Repository, hash, optionalKey)

	oururl := fmt.Sprintf("https://bitbucket.org/%v/%v", bc.Workspace, bc.Repository)

	// our descriptions are too long, so we ignore them
	body := []byte(fmt.Sprintf(`{"key":"%v", "state":"%v", "description": "%v", "url": "%v"}`, comment, st, description, oururl))

	r, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to start post: %v", err)
	}
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", bc.Token))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("Error posting github status: %v", err)
	}

	body, _ = io.ReadAll(res.Body)

	res.Body.Close()

	switch res.StatusCode {
	case 200, 201:
		// ok

	default:
		fmt.Println(res.StatusCode)
		fmt.Println(string(body))
	}
	return nil
}
