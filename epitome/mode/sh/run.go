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
)

func Run(
	cfg config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
	dynamicClient *dynamic.DynamicClient,
) error {
	// Setup input reader
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to epitomesh! Type 'help' for available commands")

	for {
		fmt.Printf("%s%s%s",
			config.ShellPromptColor,
			"epitomesh ) ",
			config.ShellResetColor)
		input, _ := reader.ReadString('\n')
		cmd := strings.TrimSpace(input)

		switch cmd {
		case "ls":
			listPods(cfg, clientset)
		case "exit":
			return nil
		case "help":
			printHelp()
		case "":
			// Do nothing for empty input
			continue
		default:
			fmt.Printf("Unknown command: %s\n", cmd)
			printHelp()
		}
	}
}

func listPods(cfg config.Config, clientset kubernetes.Interface) {
	pods, err := clientset.CoreV1().Pods(cfg.HyperdosNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		return
	}

	if len(pods.Items) == 0 {
		fmt.Println("No pods found in", cfg.HyperdosNamespace)
		return
	}

	fmt.Printf("\nPODS IN %s:\n", cfg.HyperdosNamespace)
	fmt.Printf("%-40s %-12s %-10s\n", "NAME", "STATUS", "AGE")

	for _, pod := range pods.Items {
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second)
		fmt.Printf("%-40s %-12s %-10s\n",
			pod.Name,
			getPodStatus(pod),
			age.String(),
		)
	}
	fmt.Println()
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

func printHelp() {
	fmt.Println(`
Available commands:
  ls      - List pods in current namespace
  exit    - Exit the shell
  help    - Show this help message`)
}
