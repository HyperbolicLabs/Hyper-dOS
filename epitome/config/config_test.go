package config

import (
	"net/url"
	"testing"

	env11 "github.com/caarlos0/env/v11"
	"github.com/stretchr/testify/require"
)

func TestParseEnv(t *testing.T) {
	cases := []struct {
		desc        string
		environment map[string]string
		expected    url.URL
	}{
		{
			desc: "unset",
			environment: map[string]string{
				"HYPERBOLIC_TOKEN": "token",
			},
			expected: url.URL{Scheme: "https", Host: "api.hyperbolic.xyz", Path: ""},
		},
		{
			desc: "set",
			environment: map[string]string{
				"HYPERBOLIC_GATEWAY_URL": "https://api.dev-hyperbolic.xyz",
				"HYPERBOLIC_TOKEN":       "token",
			},
			expected: url.URL{Scheme: "https", Host: "api.dev-hyperbolic.xyz", Path: ""},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			cfg := Config{}
			err := env11.ParseWithOptions(
				&cfg,
				env11.Options{
					Environment: tc.environment,
				})

			require.NoError(t, err)
			require.Equal(t, tc.expected, cfg.Jungle.HYPERBOLIC_GATEWAY_URL)
		})
	}
}
