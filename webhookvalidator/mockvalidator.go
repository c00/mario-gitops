package webhookvalidator

import (
	"io"
	"log/slog"

	"github.com/c00/mario-gitops/config"
)

var _ WebhookValidator = (*MockValidate)(nil)

type MockValidate struct {
	TagToReturn string
}

func (d *MockValidate) Validate(endpoint config.Endpoint, reader io.Reader) (string, error) {
	slog.Info("MockValidator.Validate()", "ID", endpoint.ID)

	return d.TagToReturn, nil
}
