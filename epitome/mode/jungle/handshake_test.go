package jungle

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"epitome.hyperbolic.xyz/config"
	env11 "github.com/caarlos0/env/v11"

	"github.com/stretchr/testify/require"
)

type MockHttpClient struct {
	cfg config.Config
	t   *testing.T
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.cfg.Jungle.HYPERBOLIC_TOKEN == nil {
		return nil, fmt.Errorf("HYPERBOLIC_TOKEN is not set")
	}
	require.Equal(m.t, http.MethodPost, req.Method)
	require.Equal(m.t, "https://api.hyperbolic.xyz/v1/hyperweb/login", req.URL.String())
	require.Equal(m.t, "application/json", req.Header.Get("Content-Type"))
	require.Equal(m.t, "bearer "+*m.cfg.Jungle.HYPERBOLIC_TOKEN, req.Header.Get("Authorization"))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(
			bytes.NewBufferString(
				`{"client_id":"client-id","client_secret":"client-secret","cluster_name":"cluster-name"}`)),
	}, nil
}

func Test_handshake(t *testing.T) {
	cfg := config.Config{}
	err := env11.ParseWithOptions(
		&cfg,
		env11.Options{
			Environment: map[string]string{
				"HYPERBOLIC_TOKEN": "abc123",
			},
		},
	)

	require.NoError(t, err)

	a := &agent{
		cfg: cfg,
	}

	t.Run("happy path", func(t *testing.T) {
		wantResponse := &handshakeResponse{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			ClusterName:  "cluster-name",
		}

		// expect a POST to https://api.hyperbolic.xyz/v1/hyperweb/login
		mockHttpClient := &MockHttpClient{
			cfg: cfg,
			t:   t,
		}

		gotResponse, err := a.handshake(mockHttpClient)

		require.NoError(t, err)

		require.Equal(t, wantResponse, gotResponse)
	})
}
