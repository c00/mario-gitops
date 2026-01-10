package gitops

import (
	"fmt"
	"log/slog"

	"github.com/c00/mario-gitops/gitops/github"
)

var _ GitOpser = (*GithubOps)(nil)

func NewGithubOps(token, org, repo, branch string) *GithubOps {
	return &GithubOps{
		gh: github.New(token, org, repo, branch),
	}
}

type GithubOps struct {
	gh *github.Github
}

func (g *GithubOps) Update(filepath string, jsonpath string, newTag string) error {
	slog.Info("updating github repo", "filepath", filepath, "jsonpath", jsonpath, "tag", newTag)

	sha, contents, err := g.gh.GetFileContents(filepath)
	if err != nil {
		return fmt.Errorf("cannot get file contents of '%v': %w", filepath, err)
	}

	//update the Tag
	output, err := updateYaml(contents, jsonpath, newTag)
	if err != nil {
		return fmt.Errorf("cannot update yaml: %w", err)
	}

	//commit
	_, err = g.gh.WriteFileContents(filepath, sha, fmt.Sprintf("updated yaml: %v, tag: %v", filepath, newTag), []byte(output))
	if err != nil {
		return fmt.Errorf("could not commit new file: %w", err)
	}

	return nil
}
