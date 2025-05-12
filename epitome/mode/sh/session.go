package sh

import (
	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type session struct {
	cfg           *config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	dynamicClient *dynamic.DynamicClient
	namespace     *string
	rl            *readline.Instance
	completions   readline.DynamicPrefixCompleterInterface
	cdCompletions []string
}

func (s *session) write(msg string) {
	s.rl.Write([]byte(msg))
}

func (s *session) writeln(msg string) {
	s.rl.Write([]byte(msg + "\n"))
}
