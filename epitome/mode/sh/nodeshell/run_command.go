package nodeshell

import (
	"os"
	"os/exec"
)

func RunCommand(
	sudo bool,
	args []string,
	stdin *os.File,
	stdout *os.File,
	stderr *os.File,
) error {
	if sudo {
		args = append([]string{
			"sudo",
			"-S", // -S handles pwd from stdin better
		}, args...)
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}

func RunCommandFromStr(
	sudo bool,
	command string,
	stdin *os.File,
	stdout *os.File,
	stderr *os.File,
) error {
	if sudo {
		command = "sudo" + " -S " + command
	}

	cmd := exec.Command("bash", "-c", command)

	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	return cmd.Run()
}
