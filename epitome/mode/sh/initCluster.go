package sh

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"epitome.hyperbolic.xyz/config"
	"epitome.hyperbolic.xyz/mode/sh/microk8s"
)

func (s *session) initCluster(args ...string) error {
	// use flag.parse to parse args for the -mode=cricket flag

	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	versionArg := flags.String("version", "", "Specify the version of hyperdos to install")
	roleArg := flags.String("role", "", "Specify the role to initialize the cluster in (buffalo | cricket)")
	flags.Parse(args)

	if *versionArg == "" {
		if s.cfg.HYPERDOS_VERSION == nil {
			return fmt.Errorf("must specify -version=VERSION or set $HYPERDOS_VERSION")
		} else {
			versionArg = s.cfg.HYPERDOS_VERSION
		}
	}

	if *roleArg == "" {
		return fmt.Errorf("must specify -role=(buffalo | cricket)")
	}

	roles := config.JungleRole{
		// TODO use 'contains' or split on commas so we
		// can init with multiple roles at once
		Buffalo: *roleArg == "buffalo",
		Cricket: *roleArg == "cricket",
	}

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

	// note that they underlying function should do a snap refresh --hold
	if err := s.checkAndInstallSnap("microk8s", "--classic", "--channel=1.32/stable"); err != nil {
		return err
	}

	s.write("cluster initialized\n")

	err := s.checkAndInstallHyperdos(roles, *versionArg)
	if err != nil {
		return fmt.Errorf("failed to install hyperdos: %v", err)
	}

	return nil
}

func (s *session) checkAndInstallHyperdos(roles config.JungleRole, version string) error {
	// TODO semver parse version

	s.writeln("would you like to install hyperdos now?")
	if !s.confirm() {
		return fmt.Errorf("operation canceled by user")
	}

	// note: something isn't quite smooth about s.rl.Stdout()
	microk8s.ConfigureNodeBasics(s.rl)

	// since a single baron can hold multiple jungle roles at once,
	// we check each role separately
	if s.cfg.Role.Buffalo || roles.Buffalo {
		// TODO install microceph
		return fmt.Errorf("buffalo install not yet implemented")
	}

	if s.cfg.Role.Cow || roles.Cow {
		return fmt.Errorf("cow install not yet implemented")
	}

	if s.cfg.Role.Squirrel || roles.Squirrel {
		// TODO install microceph
		return fmt.Errorf("squirrel install not yet implemented")
	}

	if s.cfg.Role.Cricket || roles.Cricket {
		s.writeln(`
		the cricket role has been selected. 
		install hyperdos with jungleRole cricket?
		Note: this will modify /var/snap/microk8s/current/args/kube-apiserver
		to expand the microk8s service IP range and nodeport range.
		`)
		if !s.confirm() {
			return fmt.Errorf("cricket setup canceled by user")
		}

		// no microceph necessary
		err := microk8s.ConfigureCricketNode(s.rl)
		if err != nil {
			return fmt.Errorf("failed to configure cricket node: %v", err)
		}

		s.writeln("microk8s has been configured for cricket mode. Would you like to helm-install hyperdos now?")
		if !s.confirm() {
			return fmt.Errorf("helm install hyperdos canceled by user")
		}

		err = microk8s.InstallHyperdos(s.cfg, version)
		if err != nil {
			return fmt.Errorf("failed to install hyperdos in cricket mode")
		}
	}

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
