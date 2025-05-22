package nodeshell

import (
	"os"
	"os/exec"
	"strings"
)

func RunCommand(
	sudo bool,
	args string,
	stdin *os.File,
	stdout *os.File,
	stderr *os.File,
) error {
	splitArgs := strings.Split(args, " ")

	if sudo {
		splitArgs = append([]string{
			"sudo",
			"-S", // -S handles pwd from stdin better
		}, splitArgs...)

	}

	cmd := exec.Command("sudo", splitArgs...)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
