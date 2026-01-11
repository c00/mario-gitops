package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/c00/mario-gitops/config"
	"github.com/c00/mario-gitops/gitops"
	"github.com/c00/mario-gitops/webhookvalidator"
)

var (
	cfg config.Config
)

func main() {
	err := run()
	if err != nil {
		slog.Error("Mario failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	// Read Environment Variables
	configPath := getEnv("MARIO_CONFIG_PATH", "config.yaml")
	GitopsToken := getEnv("MARIO_GITOPS_TOKEN", "")

	// Read config
	var err error
	cfg, err = config.GetConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config failed: %w", err)
	}
	// Probably loading of secrets could be done better.
	if cfg.GitopsToken == "" {
		cfg.GitopsToken = GitopsToken
	}

	mux := http.NewServeMux()

	// Setup endpoints
	for _, e := range cfg.Endpoints {
		// Add endpoints to the mux/
		handler, err := createWebhookHandler(e)
		if err != nil {
			return fmt.Errorf("cannot create handler for webhook: %w", err)
		}

		stub, err := url.JoinPath("/webhook", e.ID)
		if err != nil {
			return fmt.Errorf("cannot create path: %w", err)
		}

		slog.Info("Adding Webhook", "stub", stub, "name", e.Name)
		mux.HandleFunc(stub, handler)
	}

	// Not found route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug("Not found", "path", r.URL.Path)
		w.WriteHeader(404)
	})

	// Start the server
	slog.Info("Mario Serving", "port", ":8888")

	err = http.ListenAndServe(":8888", mux)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func createWebhookHandler(e config.Endpoint) (http.HandlerFunc, error) {
	// Create the validator
	var validator webhookvalidator.WebhookValidator
	switch e.ValidationStrategy {
	case config.StrategyDockerHub:
		validator = &webhookvalidator.DockerHub{}
	case config.StrategyMockValidate:
		slog.Warn("Mock Validator used")
		validator = &webhookvalidator.MockValidate{TagToReturn: "demo-tag"}
	default:
		return nil, fmt.Errorf("unsupported webhook validator strategy: %v", e.ValidationStrategy)
	}

	var gitopser gitops.GitOpser
	switch e.GitopsStrategy {
	case config.StrategyGitHub:
		gitopser = gitops.NewGithubOps(cfg.GitopsToken, cfg.GitopsOrg, cfg.GitopsRepository, cfg.GitopsBranch)
	case config.StrategyMockOps:
		slog.Warn("Mock Gitopser used")
		gitopser = &gitops.MockOps{}
	default:
		return nil, fmt.Errorf("unsupported gitopser strategy: %v", e.GitopsStrategy)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("cannot read http body", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		go func() {
			// Validate
			newTag, err := validator.Validate(e, bodyBytes)
			if err != nil {
				slog.Info("invalid webhook request", "error", err)
				return
			}

			// Run Update Actions
			for _, action := range e.Actions {
				err := gitopser.Update(action.FilePath, action.JsonPath, newTag)
				if err != nil {
					slog.Error("could not update gitops", "error", err)
					return
				}
			}

		}()

		w.WriteHeader(http.StatusOK)
	}

	return handler, nil
}
