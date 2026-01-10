package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v71/github"
)

type Github struct {
	token        string
	organization string
	client       *github.Client
	Repo         string
	Branch       string
}

func New(token, org string, repo string, branch string) *Github {
	client := github.NewClient(nil).WithAuthToken(token)

	return &Github{
		Repo:         repo,
		token:        token,
		organization: org,
		client:       client,
		Branch:       branch,
	}
}

func (g *Github) GetFileContents(path string) (sha string, contents string, err error) {
	fileContent, dirContent, _, err := g.client.Repositories.GetContents(context.Background(), g.organization, g.Repo, path, &github.RepositoryContentGetOptions{Ref: g.Branch})
	if err != nil {
		return "", "", fmt.Errorf("could not get contents: %w", err)
	}

	if fileContent != nil {
		content, err := fileContent.GetContent()
		if err != nil {
			return "", "", fmt.Errorf("could not get content: %w", err)
		}
		return fileContent.GetSHA(), content, nil
	}

	if dirContent != nil {
		return "", "", fmt.Errorf("path is a folder")
	}

	//I don't think this can happen, because if both are nil, surely an error would have been thrown.
	return "", "", fmt.Errorf("unknown error")
}

func (g *Github) WriteFileContents(path, sha, commitMsg string, content []byte) (string, error) {
	branch := g.Branch
	res, _, err := g.client.Repositories.UpdateFile(context.Background(), g.organization, g.Repo, path, &github.RepositoryContentFileOptions{
		SHA:     &sha,
		Message: &commitMsg,
		Branch:  &branch,
		Content: content,
	})
	if err != nil {
		return "", fmt.Errorf("could not post contents: %w", err)
	}

	return res.GetSHA(), nil
}
