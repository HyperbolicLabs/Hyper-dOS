package microk8s

import (
	"fmt"
	"io"
	"os"
	"strings"

	"epitome.hyperbolic.xyz/config"
	"epitome.hyperbolic.xyz/mode/sh/nodeshell"
)

// ConfigureNodeBasics configures the aspects of a microk8s node
// that are the same across all roles. For example, we might increase
// the pod limit from 110 to 200
// be careful about considerations mentioned here
//
// and here
func ConfigureNodeBasics(w io.Writer) {
	w.Write([]byte(
		`
		Note: The default Pod limit per node is 110.
        For now, we will leave it up to you if you would like to change this.
		Please refer to the following links for relevant discussion:

		- https://github.com/kubernetes/kubernetes/issues/119391
		- https://github.com/kubernetes/kubernetes/issues/23349

	`))
}

// InstallHyperdos assumes microk8s is present, properly configured,
// and running
func InstallHyperdos(cfg *config.Config, version string) error {
	sudo := true
	args := "microk8s helm repo add hyperdos https://hyperboliclabs.github.io/Hyper-dOS"
	splitArgs := strings.Split(args, " ")
	if err := nodeshell.RunCommand(sudo, splitArgs, os.Stdin, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("could not add hyperdos helm repo: %v", err)

	}

	args = "microk8s helm repo update"
	splitArgs = strings.Split(args, " ")
	if err := nodeshell.RunCommand(sudo, splitArgs, os.Stdin, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("failed to update helm repos: %v", err)
	}

	shouldEnableHyperai := false
	if cfg.Role.Buffalo {
		shouldEnableHyperai = true
	}

	args = fmt.Sprintf(`microk8s helm install hyperdos \
	hyperdos/hyperdos \
	--version %s \
	--set token="%s" \
	--set cascade.jungleRole.buffalo="%v" \
	--set cascade.jungelRole.cricket="%v"\
	--set cascade.jungelRole.cow="%v"\
	--set cascade.jungelRole.squirrel="%v"\
	--set cascade.hyperai.enabled="%v" \
	`,
		version,
		cfg.Default.HYPERBOLIC_TOKEN,
		cfg.Role.Buffalo,
		cfg.Role.Cricket,
		cfg.Role.Cow,
		cfg.Role.Squirrel,
		shouldEnableHyperai,
	)

	if err := nodeshell.RunCommandFromStr(sudo, args, os.Stdin, os.Stdout, os.Stderr); err != nil {
		return fmt.Errorf("failed to install hyperdos: %v", err)
	}

	return nil
}
