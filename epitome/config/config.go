package config

import "net/url"

type Config struct {
	// I like to use UPPER_SNAKE_CASE for config that parses from the environment,
	// as it gives a bit of intuition downstream about where these values may come from
	LOG_LEVEL              string  `env:"LOG_LEVEL" envDefault:"info"`
	HYPERBOLIC_GATEWAY_URL url.URL `env:"HYPERBOLIC_GATEWAY_URL" envDefault:"https://api.hyperbolic.xyz"`
	HYPERBOLIC_TOKEN       string  `env:"HYPERBOLIC_TOKEN,required"`
	KUBECONFIG             string  `env:"KUBECONFIG" envDefault:""`
}
