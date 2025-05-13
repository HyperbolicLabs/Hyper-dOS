package sh

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (s *session) initCluster() {
	if s.clientset != nil && !s.cfg.DEBUG {
		s.write("cluster already initialized\n")
		return
	}

	if runtime.GOOS != "linux" {
		s.writeInitNotImplementedOnThisPlatform()
		return
	}

	// check if snapd is installed
	err := s.checkAndInstallTool("snap")
	if err != nil {
		s.writeErr(err.Error())
		return
	}

	err = s.checkAndInstallSnap("microk8s", "--classic", "--channel=1.32/stable")
	if err != nil {
		s.writeErr(err.Error())
		return
	}

	s.write("cluster initialized\n")
}

func (s *session) confirm() bool {
	s.rl.SetPrompt("[Y/n]: ")
	line, err := s.rl.Readline()
	if err != nil {
		s.writeErr(err.Error())
		return false
	}

	line = strings.ToLower(strings.TrimSpace(line))

	if line == "n" || line == "no" {
		return false
	}

	if line == "" || line == "y" || line == "yes" {
		return true
	}

	return s.confirm()
}

func (s *session) checkAndInstallTool(toolName string) error {
	if _, err := exec.LookPath(toolName); err != nil {
		s.writeln(toolName + " is not installed, would you like to install it now?")
		if s.confirm() {
			return s.installTool(toolName)
		}
	}

	// already installed, nothing to do
	return nil
}

func (s *session) checkAndInstallSnap(snapName string, options ...string) error {
	if _, err := exec.LookPath(snapName); err != nil {
		s.writeln("microk8s is not installed, would you like to install it now?")
		if s.confirm() {
			return s.installSnap(snapName, options...)
		}
	}

	// already installed, nothing to do
	return nil
}

func (s *session) installSnap(snapName string, options ...string) error {
	s.writeln("installing " + snapName)

	cmd := exec.Command("sudo", append(
		[]string{
			"-S", // accept password from stdin if required
			"snap",
			"install",
			snapName,
		}, options...)...)

	// cmd.Stdin = s.rl.Terminal.GetConfig().Stdin
	cmd.Stdin = os.Stdin
	cmd.Stdout = s.rl.Stdout()
	cmd.Stderr = s.rl.Stderr()

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *session) installTool(toolName string) error {
	s.writeln("installing " + toolName)

	var cmd = exec.Command("sudo", "-S")

	// use pacman if present
	if _, err := exec.LookPath("pacman"); err == nil {
		cmd.Args = append(cmd.Args, "pacman", "-S", toolName)
	}

	// use yay if present
	if _, err := exec.LookPath("yay"); err == nil {
		cmd.Args = append(cmd.Args, "yay", toolName)
	}

	// use apt if present
	if _, err := exec.LookPath("apt"); err == nil {
		cmd = exec.Command("sudo", "apt", "install", toolName)
		cmd.Args = append(cmd.Args, "apt", "install", "-y", toolName)
	}

	if cmd == nil {
		return fmt.Errorf("could not install %s - no known package managers present", toolName)
	}

	cmd.Stdin = s.rl.Config.Stdin
	cmd.Stdout = s.rl.Stdout()
	cmd.Stderr = s.rl.Stderr()

	return cmd.Run()
}

func (s *session) writeInitNotImplementedOnThisPlatform() {
	s.writeln("******************************")
	s.writeln("This command can only be used on linux")
	s.writeln("please set up a cluster manually and restart epitomesh")
	s.writeln("")
	s.writeln("visit https://microk8s.io/tutorials for instructions to get started on macOS/Windows")
	s.writeln("******************************")
}
