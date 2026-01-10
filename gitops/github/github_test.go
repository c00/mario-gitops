//go:build integration

package github

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/c00/mario-gitops/config"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestGithub_GetFileContents(t *testing.T) {
	loadSettings()

	cfg := config.Config{
		GitopsToken:      os.Getenv("MARIO_GITOPS_TOKEN"),
		GitopsOrg:        os.Getenv("MARIO_GITOPS_ORG"),
		GitopsRepository: os.Getenv("MARIO_GITOPS_REPO"),
		GitopsBranch:     os.Getenv("MARIO_GITOPS_BRANCH"),
	}

	gh := New(cfg.GitopsToken, cfg.GitopsOrg, cfg.GitopsRepository, cfg.GitopsBranch)

	expected := `This is a test file for reading.`

	sha, contents, err := gh.GetFileContents("testdata/read-file.txt")
	assert.Nil(t, err)
	assert.Equal(t, expected, contents)
	assert.NotEqual(t, "", sha)
}

func TestGithub_WriteFileContents(t *testing.T) {
	loadSettings()

	cfg := config.Config{
		GitopsToken:      os.Getenv("MARIO_GITOPS_TOKEN"),
		GitopsOrg:        os.Getenv("MARIO_GITOPS_ORG"),
		GitopsRepository: os.Getenv("MARIO_GITOPS_REPO"),
		GitopsBranch:     os.Getenv("MARIO_GITOPS_BRANCH"),
	}

	gh := New(cfg.GitopsToken, cfg.GitopsOrg, cfg.GitopsRepository, cfg.GitopsBranch)

	path := "testdata/write-file.txt"

	sha, contents, err := gh.GetFileContents(path)
	assert.Nil(t, err)
	assert.NotEmpty(t, sha)

	contents = fmt.Sprintf("This file was last written: %v", time.Now().String())
	commit, err := gh.WriteFileContents(path, sha, "TestGithub_WriteFileContents: test run", []byte(contents))
	assert.Nil(t, err)
	assert.NotEmpty(t, commit)
}

func loadSettings() {
	folder, err := findProjectRoot()
	if err != nil {
		panic(err)
	}

	godotenv.Load(filepath.Join(folder, "test.env"))
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
