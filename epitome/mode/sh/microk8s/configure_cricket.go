package microk8s

import (
	"fmt"
	"io"
	"os"

	"epitome.hyperbolic.xyz/mode/sh/nodeshell"
)

func ConfigureCricketNode(w io.Writer) error {
	sudo := true

	// 1) edit /var/snap/microk8s/current/args/kube-apiserver
	// to change the service IP range to /16 instead of /24
	if err := upgradeServiceIPRange(
		sudo,
		"/var/snap/microk8s/current/args/kube-apiserver",
	); err != nil {
		return fmt.Errorf("failed to upgrade service IP range: %v", err)
	}

	// Note: ideally, we would be able to specify multiple ranges,
	// but kube-apiserver doesn't support that. could be low-hanging fruit for a good PR
	// e.g. 80,443,30000-40000
	if err := upgradeNodePortRange(
		sudo,
		"/var/snap/microk8s/current/args/kube-apiserver",
		"80-50000",
	); err != nil {
		return fmt.Errorf("failed to upgrade node port range: %v", err)
	}

	// 4) microk8s refresh-certs --cert server.crt

	// 5) microk8s stop && microk8s start
	if err := nodeshell.RunCommandFromStr(
		sudo,
		"microk8s stop && microk8s start",
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		return fmt.Errorf("failed to restart microk8s: %v", err)
	}

	return nil
}

func buildServiceIPRangeSedCommand(path string) []string {
	return []string{
		"sed",
		"-i",
		"s/^--service-cluster-ip-range=\\([0-9]\\+\\)\\.\\([0-9]\\+\\)\\.\\([0-9]\\+\\)\\.\\([0-9]\\+\\)\\/24$/--service-cluster-ip-range=\\1\\.\\2\\.0\\.0\\/16/",
		path,
	}
}

func upgradeServiceIPRange(sudo bool, path string) error {
	command := buildServiceIPRangeSedCommand(path)
	return nodeshell.RunCommand(
		sudo,
		command,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
}

func upgradeNodePortRange(sudo bool, path string, newRange string) error {

	// just delete the old one and append the new one
	// https://stackoverflow.com/a/29314671
	err := nodeshell.RunCommand(
		sudo,
		[]string{
			"sed",
			"-i",
			`/--service-node-port-range=/d`,
			path,
		},
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	if err != nil {
		return err
	}

	// add a newline if necessary
	command := "sed -i '$a\\' " + path
	err = nodeshell.RunCommandFromStr(
		sudo,
		command,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
	if err != nil {
		return err
	}

	command = "echo '--service-node-port-range=" + newRange + "' >> " + path

	err = nodeshell.RunCommandFromStr(
		sudo,
		command,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	return err
}
