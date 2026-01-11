package webhookvalidator

import (
	"log/slog"

	"github.com/c00/mario-gitops/config"
)

var _ WebhookValidator = (*MockValidate)(nil)

type MockValidate struct {
	TagToReturn string
}

func (d *MockValidate) Validate(endpoint config.Endpoint, reader []byte) (string, error) {
	slog.Info("MockValidator.Validate()", "ID", endpoint.ID)

	return d.TagToReturn, nil
}
