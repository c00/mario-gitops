//go:build integration

package github

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/c00/mario-gitops/config"
	"github.com/stretchr/testify/assert"
)

func TestGithub_GetFileContents(t *testing.T) {
	//todo dotenv test stuff.

	cfg := config.Config{
		GitopsToken:      os.Getenv("MARIO_GITOPS_TOKEN"),
		GitopsOrg:        os.Getenv("MARIO_GITOPS_ORG"),
		GitopsRepository: os.Getenv("MARIO_GITOPS_REPO"),
		GitopsBranch:     os.Getenv("MARIO_GITOPS_BRANCH"),
	}

	gh := New(cfg.GitopsToken, cfg.GitopsOrg, cfg.GitopsRepository, cfg.GitopsBranch)

	expected := `foo: bar
wow:
  - such
  - great
  - file
`

	sha, contents, err := gh.GetFileContents("src/read-test-file.yaml")
	assert.Nil(t, err)
	assert.Equal(t, expected, contents)
	assert.NotEqual(t, "", sha)
}

func TestGithub_WriteFileContents(t *testing.T) {
	//todo dotenv test stuff.

	cfg := config.Config{
		GitopsToken:      os.Getenv("MARIO_GITOPS_TOKEN"),
		GitopsOrg:        os.Getenv("MARIO_GITOPS_ORG"),
		GitopsRepository: os.Getenv("MARIO_GITOPS_REPO"),
		GitopsBranch:     os.Getenv("MARIO_GITOPS_BRANCH"),
	}

	gh := New(cfg.GitopsToken, cfg.GitopsOrg, cfg.GitopsRepository, cfg.GitopsBranch)

	path := "src/write-test-file.yaml"

	sha, contents, err := gh.GetFileContents(path)
	assert.Nil(t, err)
	assert.NotEmpty(t, sha)

	contents = fmt.Sprintf("%v\n%v", contents, time.Now().String())
	commit, err := gh.WriteFileContents(path, sha, "TestGithub_WriteFileContents: test run", []byte(contents))
	assert.Nil(t, err)
	assert.NotEmpty(t, commit)
}
