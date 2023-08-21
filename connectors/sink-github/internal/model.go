package internal

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type IssueData struct {
	URL  string `json:"url"`
	Body string `json:"body"`
}

func (m IssueData) GetIssue() (*Issue, error) {
	// https://api.github.com/repos/octocat/Hello-World/issues/1
	url, err := url.Parse(m.URL)
	if err != nil {
		return nil, err
	}
	paths := strings.Split(url.Path, "/")
	if len(paths) != 6 {
		return nil, fmt.Errorf("url is invalid")
	}
	num, err := strconv.Atoi(strings.TrimSpace(paths[5]))
	if err != nil {
		return nil, fmt.Errorf("url issue nubmer is invalid")
	}
	return &Issue{
		Owner:  paths[2],
		Repo:   paths[3],
		Number: num,
	}, nil
}

type Issue struct {
	Owner  string
	Repo   string
	Number int
}
