package webhookvalidator

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/c00/mario-gitops/config"
)

var _ WebhookValidator = (*DockerHub)(nil)

type DockerHub struct {
}

func (d *DockerHub) Validate(endpoint config.Endpoint, rawBody []byte) (string, error) {
	body := DockerHookPayload{}

	err := json.Unmarshal(rawBody, &body)
	if err != nil {
		return "", fmt.Errorf("cannot unmarshall body: %w", err)
	}

	slog.Info("Docker Webhook Received", "callback URL", body.CallbackURL)

	// latest tags should be ignored
	if body.PushData.Tag == "latest" {
		return "", fmt.Errorf("dont trigger on 'latest' tag")
	}

	// Check repo name
	if endpoint.DockerRepository != body.Repository.RepoName {
		return "", fmt.Errorf("configure repo and body repo do not match: %v != %v", endpoint.DockerRepository, body.Repository.RepoName)
	}

	// TODO Figure out how the fuck these fucking webhooks work, and then re-enable it. Until then, just yolo it
	// Check webhook validity
	// response, err := http.Post(body.CallbackURL, "application/json", bytes.NewBufferString(`{"status":"success"}`))
	// if err != nil {
	// 	return "", fmt.Errorf("docker callback failed: %w", err)
	// }

	// if response.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("not legitimate docker webhook. Got status code %v", response.StatusCode)
	// }

	return body.PushData.Tag, nil
}

type DockerHookPayload struct {
	CallbackURL string     `json:"callback_url"`
	PushData    PushData   `json:"push_data"`
	Repository  Repository `json:"repository"`
}

type PushData struct {
	PushedAt int64  `json:"pushed_at"`
	Pusher   string `json:"pusher"`
	Tag      string `json:"tag"`
}

type Repository struct {
	CommentCount    int    `json:"comment_count"`
	DateCreated     int64  `json:"date_created"`
	Description     string `json:"description"`
	Dockerfile      string `json:"dockerfile"`
	FullDescription string `json:"full_description"`
	IsOfficial      bool   `json:"is_official"`
	IsPrivate       bool   `json:"is_private"`
	IsTrusted       bool   `json:"is_trusted"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	Owner           string `json:"owner"`
	RepoName        string `json:"repo_name"`
	RepoURL         string `json:"repo_url"`
	StarCount       int    `json:"star_count"`
	Status          string `json:"status"`
}
