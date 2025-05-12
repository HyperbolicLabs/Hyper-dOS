package sh

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
)

func (s *session) initReadline() error {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("ls"),
		readline.PcItem("cd"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("clear"),
	)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          s.getPrompt(),
		AutoComplete:    completer,
		HistoryFile:     "/tmp/epitome_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	s.rl = l
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

func (s *session) listPods(clientset kubernetes.Interface, namespace string) {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		return
	}

	if len(pods.Items) == 0 {
		s.write(fmt.Sprintf("No pods found in %s\n", namespace))
		return
	}

	s.write(fmt.Sprintf("\nPODS IN %s:\n", namespace))
	s.write(fmt.Sprintf("%-40s %-12s %-10s\n", "NAME", "STATUS", "AGE"))

	for _, pod := range pods.Items {
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second)
		s.write(fmt.Sprintf("%-40s %-12s %-10s\n",
			pod.Name,
			getPodStatus(pod),
			age.String(),
		))
	}

	s.write("\n")
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
