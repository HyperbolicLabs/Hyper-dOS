package sh

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (s *session) initCluster(args ...string) error {
	// use flag.parse to parse args for the -mode=cricket flag

	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	modeArg := flags.String("mode", "ronin", "Specify the mode to initialize the cluster in (ronin | buffalo | cricket)")
	flags.Parse(args)

	if s.clientset != nil && !s.cfg.DEBUG {
		s.write("cluster already initialized\n")
		return nil
	}

	if runtime.GOOS != "linux" {
		s.writeInitNotImplementedOnThisPlatform()
		return fmt.Errorf("not implemented on this platform")
	}

	// check if snapd is installed
	if err := s.checkAndInstallTool("snap"); err != nil {
		return err
	}

	if err := s.checkAndInstallSnap("microk8s", "--classic", "--channel=1.32/stable"); err != nil {
		return err
	}

	// switch modearg
	switch *modeArg {
	case "ronin":
		return fmt.Errorf("ronin mode not yet implemented")
	case "buffalo":
		return fmt.Errorf("buffalo mode not yet implemented")
	case "cricket":
		return s.installHyperdos(
			s.cfg.HyperdosNamespace,
			"TODO cricket proto",
		)
	}

	s.write("cluster initialized\n")

	return nil
}

func (s *session) installHyperdos(
	namespace string,
	role string) error {

	return fmt.Errorf("TODO: cricket install not implemented")
}

func (s *session) confirm() bool {
	s.rl.SetPrompt("[Y/n]: ")
	line, err := s.rl.Readline()
	if err != nil {
		s.writeErr(err)
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
	hold := true // for now, hold all snaps
	if _, err := exec.LookPath(snapName); err != nil {
		s.writeln("microk8s is not installed, would you like to install it now?")
		if s.confirm() {
			return s.installSnap(snapName, hold, options...)
		}
	}

	// already installed, nothing to do
	return nil
}

func (s *session) installSnap(
	snapName string,
	hold bool,
	options ...string) error {
	s.writeln("installing " + snapName)

	cmd := exec.Command("sudo", append(
		[]string{
			"-S", // accept password from stdin if required
			"snap",
			"install",
			snapName,
		}, options...)...)

	// for some reason s.rl.Terminal.GetConfig().Stdin doesn't work smoothly
	// nor do s.rl.Stdout() or s.rl.Stderr()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	} else if hold {
		// if requested, hold the snap on a successful install
		cmd = exec.Command("sudo",
			"snap", "refresh", "--hold", snapName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
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

	// for some reason s.rl.Terminal.GetConfig().Stdin doesn't work smoothly
	// nor do s.rl.Stdout() or s.rl.Stderr()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
