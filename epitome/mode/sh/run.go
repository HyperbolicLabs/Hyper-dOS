package sh

import (
	"io"
	"strings"

	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
) error {

	s, closeFunc, err := NewSession(
		&cfg,
		logger,
		clientset)
	if err != nil {
		return err
	}
	defer closeFunc()

	s.rl.Write([]byte("Welcome to epitomesh! Type 'help' for available commands\n"))

	for {
		line, err := s.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)

		// parse into command and args by splitting on spaces
		parts := strings.Split(line, " ")
		cmd := strings.ToLower(parts[0])

		var args []string
		if len(parts) > 1 {
			args = parts[1:]
		}

		err = s.dispatch(cmd, args...)
		if err != nil {
			s.writeErr(err)
		}

		// Update prompt after command execution
		s.rl.SetPrompt(s.getPrompt())
	}

	return nil
}

func getPodStatus(pod corev1.Pod) string {
	if pod.DeletionTimestamp != nil {
		return "Terminating"
	}
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return "Ready"
		}
	}
	return string(pod.Status.Phase)
}
func (s *session) printHelp() {
	s.writeln(`
Available commands:
  init      - Initialize a new hyperdos cluster 
  clear     - Clear the screen
  ls        - List resources in current context
  cd <ns>   - Enter a namespace context
  cd ..     - Return to namespace list view
  exit      - Exit the shell
  help      - Show this help message
`)
}
