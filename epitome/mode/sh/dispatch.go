package sh

import (
	"fmt"

	"github.com/chzyer/readline"
)

func (s *session) dispatch(cmd string, args ...string) error {
	switch cmd {
	case "init":
		if err := s.initCluster(args...); err != nil {
			return err
		}
	case "cd":
		if err := s.cd(args...); err != nil {
			return err
		}
	case "ls":
		if err := s.ls(args...); err != nil {
			return err
		}
	case "clear":
		readline.ClearScreen(s.rl)
	case "exit":
		return nil
	case "help":
		s.printHelp()
	case "":
		// Do nothing for empty input
		// continue
	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}

	return nil
}
