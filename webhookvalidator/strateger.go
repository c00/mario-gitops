package webhookvalidator

import (
	"github.com/c00/mario-gitops/config"
)

// WebhookValidator takes the uuid and payload and verifies that the webhook request was legitimate
type WebhookValidator interface {
	Validate(endpoint config.Endpoint, payload []byte) (string, error)
}
