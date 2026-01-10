package mockhttp

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockHttp(t *testing.T) {
	m := MockHttp{Port: 9999}

	m.Start()
	defer m.Stop()

	res, err := http.Get("http://localhost:9999/ok")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get("http://localhost:9999/notfound")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	res, err = http.Get("http://localhost:9999/error")
	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
}
