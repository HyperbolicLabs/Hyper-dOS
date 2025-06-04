package config

import (
	"net/url"
	"time"
)

const (
	CPUNameLabelKey  = "hyperbolic.xyz/cpu-name"
	CalicoNamespace  = "kube-system"
	ShellPromptColor = "\033[32m" // Green text
	ShellResetColor  = "\033[0m"  // Reset colors
)

type Config struct {
	Default  DefaultMode
	Monkey   MonkeyMode
	Maintain MaintainMode
	Shell    ShellMode
	Role     JungleRole
	// I like to use UPPER_SNAKE_CASE for config that parses from the environment,
	// as it gives a bit of intuition downstream about where these values may come from
	LOG_LEVEL         string  `env:"LOG_LEVEL" envDefault:"info"`
	DEBUG             bool    `env:"DEBUG" envDefault:"false"`
	KUBECONFIG        string  `env:"KUBECONFIG" envDefault:""`
	HyperwebNamespace string  `env:"HYPERWEB_NAMESPACE" envDefault:"hyperweb"`
	HyperdosNamespace string  `env:"HYPERDOS_NAMESPACE" envDefault:"hyperdos"`
	HYPERDOS_VERSION  *string `env:"HYPERDOS_VERSION"`
}

type JungleRole struct {
	Buffalo  bool `env:"JUNGLE_ROLE_BUFFALO" envDefault:"false"`
	Cow      bool `env:"JUNGLE_ROLE_COW" envDefault:"false"`
	Cricket  bool `env:"JUNGLE_ROLE_CRICKET" envDefault:"false"`
	Squirrel bool `env:"JUNGLE_ROLE_SQUIRREL" envDefault:"false"`
}

type DefaultMode struct {
	ReconcileInterval      time.Duration `env:"DEFAULT_RECONCILE_INTERVAL" envDefault:"1m"`
	HYPERBOLIC_GATEWAY_URL url.URL       `env:"HYPERBOLIC_GATEWAY_URL" envDefault:"https://api.hyperbolic.xyz"`
	HYPERBOLIC_TOKEN       *string       `env:"HYPERBOLIC_TOKEN"`
}

type MaintainMode struct {
	ReconcileInterval time.Duration `env:"MAINTAIN_RECONCILE_INTERVAL" envDefault:"24h"`
}

type MonkeyMode struct {
	ReconcileInterval    time.Duration `env:"MONKEY_RECONCILE_INTERVAL" envDefault:"1m"`
	KUBERNETES_NODE_NAME string        `env:"KUBERNETES_NODE_NAME,required" envDefault:""`
}

type ShellMode struct {
	VIM_MODE bool `env:"VIM_MODE" envDefault:"false"`
}
