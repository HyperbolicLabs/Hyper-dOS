package microk8s

import (
	"fmt"
	"io"
)

func ConfigureCricketNode(w io.Writer) error {
	// 1) edit /var/snap/microk8s/current/args/kube-apiserver
	// to change the service IP range to /16 instead of /24

	// 2) edit /var/snap/microk8s/current/args/kube-apiserver
	// to change the nodeport range to 80,443,30000-40000

	// 3) TODO edit /var/snap/microk8s/current/args/kube-apiserver
	// to change the max node count from 110 to 109 (just for now, until we figure out what the actual count we want is)

	// 4) microk8s refresh-certs --cert server.crt

	// 5) microk8s stop && microk8s start
	return fmt.Errorf("not implemented")
}
