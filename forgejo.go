package main

import (
	"fmt"
	"net/http"
	"bytes"
	"io"
)

type ForgejoConfiguration struct {
	Domain string
	User string
	Repository string
	Token string
	Ssh bool
}

func (fc *ForgejoConfiguration) Valid() bool {
	if fc == nil {
		return false
	}
	return fc.Domain != "" && fc.User != "" && fc.Repository != ""
}

func (fc *ForgejoConfiguration)	SetStatus(status CiStatus, comment, description, hash string) error {
	if fc.Token == "" {
		return nil
	}
	url := fmt.Sprintf("https://%v/api/v1/repos/%v/%v/statuses/%v", fc.Domain, fc.User, fc.Repository, hash)

	forgejoStatus := "error"
	switch status {
		case KInProgress:
			forgejoStatus = "pending"
		case KFailed:
			forgejoStatus = "failure"
		case KError:
			forgejoStatus = "error"
		case KSucceeded:
			forgejoStatus = "success"
	}

	if len(description) > 255 {
		description = description[:253] + "..."
	}

	body := []byte(fmt.Sprintf(`{"state":"%s","context": "%s","description":"%s"}`, forgejoStatus, comment, description))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Failed to start post: %v", err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("token %s", fc.Token))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("Error posting forgejo status: %v", err)
	}

	body, _ = io.ReadAll(res.Body)

	res.Body.Close()

	if res.StatusCode != 200 && res.StatusCode != 201 {
		fmt.Println(res.StatusCode)
		fmt.Println(string(body))
	}
	return nil	
}

func (fc *ForgejoConfiguration)	NixUrl(revision string) string {
	if !fc.Ssh {
		return fmt.Sprintf("git+https://%s/%s/%s?rev=%s", fc.Domain, fc.User, fc.Repository, revision)
	}
	return fmt.Sprintf("git+ssh://git@%s/%s/%s?rev=%s", fc.Domain, fc.User, fc.Repository, revision)
}

func (fc *ForgejoConfiguration) GitUrl() string {
	if !fc.Ssh {
		return fmt.Sprintf("https://%s/%s/%s.git", fc.Domain, fc.User, fc.Repository)
	}
	return fmt.Sprintf("git@%s:%s/%s.git", fc.Domain, fc.User, fc.Repository)
}
