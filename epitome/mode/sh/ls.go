package sh

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func (s *session) ls(target *string) {
	if target != nil {
		s.writeln("TODO: ls <target> not implemented")
		return
	}

	// if we are not in a namespace, list namespaces
	if s.namespace == nil {
		s.listNamespaces(s.clientset)
	} else {
		s.listPods(s.clientset, *s.namespace)
	}

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

		s.cdCompletions = append(s.cdCompletions, ns.Name)
	}

	s.write("\n")
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

		s.cdCompletions = append(s.cdCompletions, pod.Name)
	}

	// Refresh readline with new completions
	s.rl.Config.AutoComplete = s.completions
	s.rl.Refresh()

	s.write("\n")
}
