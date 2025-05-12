package sh

import (
	"fmt"
	"io"
	"strings"

	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
)

func (s *session) initReadline() error {
	// Initialize with static commands first
	s.completions = readline.NewPrefixCompleter(
		readline.PcItem("cd",
			readline.PcItemDynamic(s.getCdCompletions),
		),
		readline.PcItem("ls"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("clear"),
	)

	readlineInstance, err := readline.NewEx(&readline.Config{
		Prompt:          s.getPrompt(),
		AutoComplete:    s.completions, // Use our dynamic completer
		HistoryFile:     "/tmp/epitome_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	s.rl = readlineInstance
	return nil
}

func (s *session) getPrompt() string {
	path := ""
	if s.namespace != nil {
		path = "/" + *s.namespace
	}
	return fmt.Sprintf("%sepitomesh%s%s )%s ",
		config.ShellPromptColor,
		path,
		config.ShellResetColor,
		config.ShellResetColor)
}

// TODO migrate to readline?
// https://github.com/chzyer/readline/blob/main/example/readline-demo/readline-demo.go

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient *dynamic.DynamicClient,
) error {

	s := &session{
		cfg:           &cfg,
		logger:        logger,
		clientset:     clientset,
		dynamicClient: dynamicClient,
	}

	if err := s.initReadline(); err != nil {
		return err
	}
	defer s.rl.Close()

	s.rl.Write([]byte("Welcome to epitomesh! Type 'help' for available commands\n"))

	for {
		line, err := s.rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			break
		}

		cmd := strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(cmd, "cd "):
			parts := strings.SplitN(cmd, " ", 2)
			if len(parts) != 2 {
				fmt.Println("Usage: cd <namespace>")
				continue
			}

			s.cd(&parts[1])
		case cmd == "ls":
			s.ls(nil)
		case strings.HasPrefix(cmd, "ls "):
			parts := strings.SplitN(cmd, " ", 2)
			if len(parts) != 2 {
				s.writeln("Usage: ls <target>")
				continue
			}
			// user passed an argument to ls, we should use it
			s.ls(ptr.To(parts[1]))
		case cmd == "clear":
			readline.ClearScreen(s.rl)
		case cmd == "exit":
			return nil
		case cmd == "help":
			s.printHelp()
		case cmd == "":
			// Do nothing for empty input
			continue
		default:
			s.writeln(fmt.Sprintf("Unknown command: %s", cmd))
			s.printHelp()
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
  clear     - Clear the screen
  ls        - List resources in current context
  cd <ns>   - Enter a namespace context
  cd ..     - Return to namespace list view
  exit      - Exit the shell
  help      - Show this help message
`)
}
