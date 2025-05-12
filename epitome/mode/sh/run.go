package sh

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/ptr"
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient *dynamic.DynamicClient,
) error {

	// Setup input reader
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	s := session{
		cfg:           &cfg,
		logger:        logger,
		clientset:     clientset,
		dynamicClient: dynamicClient,
		namespace:     nil,
		reader:        reader,
		writer:        writer,
	}

	fmt.Println("Welcome to epitomesh! Type 'help' for available commands")

	for {
		s.prompt()
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

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
			fmt.Print("\033[H\033[2J") // ANSI escape code to clear screen
		case cmd == "exit":
			return nil
		case cmd == "help":
			s.printHelp()
		case cmd == "":
			// Do nothing for empty input
			continue
		default:
			s.write(fmt.Sprintf("Unknown command: %s\n", cmd))
			s.printHelp()
		}
	}
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
	s.writer.WriteString(`
Available commands:
  clear     - Clear the screen
  ls        - List resources in current context
  cd <ns>   - Enter a namespace context
  cd ..     - Return to namespace list view
  exit      - Exit the shell
  help      - Show this help message
`)
}
