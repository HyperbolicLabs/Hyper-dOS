package sh

import (
	"bufio"
	"context"
	"fmt"
	"time"

	"epitome.hyperbolic.xyz/config"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type session struct {
	cfg           *config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	dynamicClient *dynamic.DynamicClient
	namespace     *string
	reader        *bufio.Reader
	writer        *bufio.Writer
}

func (s *session) path() string {
	if s.namespace == nil {
		return ""
	}
	return "/" + *s.namespace
}

func (s *session) prompt() {
	s.write(fmt.Sprintf("%s%s%s%s ) ",
		config.ShellPromptColor,
		"epitomesh",
		s.path(),
		config.ShellResetColor))
}

func (s *session) write(msg string) {
	s.writer.WriteString(msg)
	s.writer.Flush()
}

func (s *session) writeln(msg string) {
	s.write(msg + "\n")
}

func (s *session) listNamespaces(clientset kubernetes.Interface) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing namespaces: %v\n", err)
		return
	}

	s.write("\nNAMESPACES:\n")
	s.write(fmt.Sprintf("%-30s %-12s %-10s\n", "NAME", "STATUS", "AGE"))

	for _, ns := range namespaces.Items {
		age := time.Since(ns.CreationTimestamp.Time).Round(time.Second)
		status := string(ns.Status.Phase)
		if ns.DeletionTimestamp != nil {
			status = "Terminating"
		}

		s.write(fmt.Sprintf("%-30s %-12s %-10s\n",
			ns.Name,
			status,
			age.String(),
		))
	}

	s.write("\n")
}
