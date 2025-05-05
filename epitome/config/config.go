package config

import (
	"net/url"
	"time"
)

const (
	CPUNameLabelKey = "hyperbolic.xyz/cpu-name"
)

type Config struct {
	Default  DefaultMode
	Monkey   MonkeyMode
	Maintain MaintainMode
	// I like to use UPPER_SNAKE_CASE for config that parses from the environment,
	// as it gives a bit of intuition downstream about where these values may come from
	LOG_LEVEL         string `env:"LOG_LEVEL" envDefault:"info"`
	KUBECONFIG        string `env:"KUBECONFIG" envDefault:""`
	HyperwebNamespace string `env:"HYPERWEB_NAMESPACE" envDefault:"hyperweb"`
	HyperdosNamespace string `env:"HYPERDOS_NAMESPACE" envDefault:"hyperdos"`
}

type DefaultMode struct {
	ReconcileInterval      time.Duration `env:"DEFAULT_RECONCILE_INTERVAL" envDefault:"1m"`
	HYPERBOLIC_GATEWAY_URL url.URL       `env:"HYPERBOLIC_GATEWAY_URL" envDefault:"https://api.hyperbolic.xyz"`
	HYPERBOLIC_TOKEN       string        `env:"HYPERBOLIC_TOKEN,required"`
}

type MaintainMode struct {
	ReconcileInterval time.Duration `env:"MAINTAIN_RECONCILE_INTERVAL" envDefault:"24h"`
}

type MonkeyMode struct {
	ReconcileInterval    time.Duration `env:"MONKEY_RECONCILE_INTERVAL" envDefault:"1m"`
	KUBERNETES_NODE_NAME string        `env:"KUBERNETES_NODE_NAME,required" envDefault:""`
}
