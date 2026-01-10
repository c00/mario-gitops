package config

import (
	"errors"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

const (
	// Validators
	// StrategyDockerHub expects a webhook from Docker Hub
	StrategyDockerHub    = "dockerhub"
	StrategyMockValidate = "mock"

	// Gitopsers
	StrategyGitHub  = "github"
	StrategyMockOps = "mock"
)

type Config struct {
	GitopsOrg        string     `yaml:"gitopsOrg"`
	GitopsRepository string     `yaml:"gitopsRepository"`
	GitopsBranch     string     `yaml:"gitopsBranch"`
	GitopsToken      string     `yaml:"gitopsToken"`
	Endpoints        []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	// Name is a user-friendly name of the endpoint. Used in logs.
	Name string `yaml:"name"`
	// ID is a unique ID used in the url.
	ID string `yaml:"id"`
	// ValidationStrategy determines what to do when a webhook is received. dockerhub and mock is supported
	ValidationStrategy string `yaml:"validationStrategy"`
	// DockerRepository indicates which container repository is expected.
	DockerRepository string `yaml:"dockerRepository"`

	// GitopsStrategy indicates how the actions will be taken. github and mock are supported
	GitopsStrategy string `yaml:"gitopsStrategy"`

	// Actions to take in gitops repo
	Actions []GitopsAction `yaml:"actions"`
}

type GitopsAction struct {
	// The path to the yaml file in the gitops repo
	FilePath string `yaml:"filePath"`
	// The JsonPath query to the image tag
	JsonPath string `yaml:"jsonPath"`
}

func (c Config) Validate() error {
	if c.GitopsRepository == "" {
		return fmt.Errorf("gitopsRespository cannot be empty")
	}

	if len(c.Endpoints) == 0 {
		return fmt.Errorf("no endpoints configured")
	}

	for _, e := range c.Endpoints {
		err := e.Validate()
		if err != nil {
			return fmt.Errorf("invalid endpoint: %w", err)
		}
	}

	return nil
}

func (e Endpoint) Validate() error {
	if e.Name == "" {
		return errors.New("endpoint name not set")
	}

	if e.ID == "" {
		return fmt.Errorf("endpoint ID not set for endpoint: %v", e.Name)
	}

	if e.GitopsStrategy == "" {
		return fmt.Errorf("gitopsStategy not set for endpoint: %v", e.Name)
	}

	if e.ValidationStrategy == "" {
		return fmt.Errorf("endpoint strategy not set for endpoint %v", e.Name)
	}

	if e.DockerRepository == "" {
		return fmt.Errorf("docker repository cannot be empty for endpoint %v", e.Name)
	}

	if len(e.Actions) == 0 {
		return fmt.Errorf("no actions defined for endpoint %v", e.Name)
	}

	for _, a := range e.Actions {
		err := a.Validate()
		if err != nil {
			return fmt.Errorf("action not valid for endpoint '%v': %w", e.Name, err)
		}
	}

	return nil
}

func (e GitopsAction) Validate() error {
	if e.FilePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}

	if e.JsonPath == "" {
		return fmt.Errorf("json path cannot be empty")
	}

	return nil
}

func GetConfig(filepath string) (Config, error) {
	if filepath == "" {
		return Config{}, errors.New("config filepath (MARIO_CONFIG_PATH) cannot be empty")
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, fmt.Errorf("cannot read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("cannot unmarshall config: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return Config{}, fmt.Errorf("config invalid: %w", err)
	}

	return config, nil
}
