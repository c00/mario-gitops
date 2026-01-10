package webhookvalidator_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/c00/mario-gitops/config"
	"github.com/c00/mario-gitops/mockhttp"
	"github.com/c00/mario-gitops/webhookvalidator"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDockerHub_Validate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		endpoint config.Endpoint
		reader   io.Reader
		wantTag  string
		wantErr  bool
	}{
		{
			name:     "base",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9999/ok", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "v1"}}`),
		},
		{
			name:     "error on latest tag",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "latest",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9999/ok", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "latest"}}`),
			wantErr:  true,
		},
		{
			name:     "unreachable server",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9900/ok", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "v1"}}`),
			wantErr:  true,
		},
		{
			name:     "callback error",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9999/error", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "v1"}}`),
			wantErr:  true,
		},
		{
			name:     "callback not found",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9999/notfound", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "v1"}}`),
			wantErr:  true,
		},
		{
			name:     "wrong repo",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/wrongrepo"},
			wantTag:  "v1",
			reader:   bytes.NewBufferString(`{"callback_url": "http://127.0.0.1:9999/ok", "repository":	{"repo_name": "demo/repo"}, "push_data": {"tag": "v1"}}`),
			wantErr:  true,
		},
		{
			name:     "wrong body",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1.0.0",
			reader:   bytes.NewBufferString(`{"foo": "bar"}`),
			wantErr:  true,
		},
		{
			name:     "not json body",
			endpoint: config.Endpoint{ID: "1a2b3c", DockerRepository: "demo/repo"},
			wantTag:  "v1.0.0",
			reader:   bytes.NewBufferString(`not json`),
			wantErr:  true,
		},
		// Todo add test cases
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := mockhttp.MockHttp{Port: 9999}
			ms.Start()
			defer ms.Stop()

			var d webhookvalidator.DockerHub
			gotTag, err := d.Validate(tt.endpoint, tt.reader)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantTag, gotTag)
		})
	}
}
