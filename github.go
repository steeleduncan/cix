/*
   github.go - Github tools for Cix

   Copyright 2024 Duncan Steele

   Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

   The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

   THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */
package main

import (
	"io"
	"fmt"
	"bytes"
	"net/http"
)

type GithubConfiguration struct {
	User string
	Repository string
	StatusPat string
}
var _ RepoSource = &GithubConfiguration {}

func (gc *GithubConfiguration) NixUrl(revision string) string {
	return fmt.Sprintf("github:%v/%v?rev=%v", gc.User, gc.Repository, revision)
}

func (gc *GithubConfiguration) GitUrl() string {
	return fmt.Sprintf("git@github.com:%v/%v", gc.User, gc.Repository)
}

func (gc *GithubConfiguration) Valid() bool {
	if gc == nil {
		return false
	}

	return gc.User != "" && gc.Repository != ""
}

func (gc *GithubConfiguration) SetStatus(status CiStatus, comment, description, hash string) error {
	if gc.StatusPat == "" {
		return nil
	}
	url := fmt.Sprintf("https://api.github.com/repos/%v/%v/statuses/%v", gc.User, gc.Repository, hash)

	st := "error"
	switch status {
	case KError:
		st = "error"

	case KInProgress:
		st = "pending"

	case KFailed:
		st = "failure"

	case KSucceeded:
		st = "success"
	}

	if len(description) > 140 {
		description = description[:136] + "..."
	}

	// our descriptions are too long, so we ignore them
	body := []byte(fmt.Sprintf(`{"state":"%v","context":"%v", "description": "%v"}`, st, comment, description))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to start post: %v", err)
	}
	r.Header.Add("Accept", "application/vnd.github+json")
	r.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", gc.StatusPat))

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
