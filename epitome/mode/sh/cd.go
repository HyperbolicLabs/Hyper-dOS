package sh

import "fmt"

// func (s *session) cd(destination *string) {
func (s *session) cd(args ...string) error {
	if s.clientset == nil {
		return fmt.Errorf("no cluster connected")
	}

	if len(args) > 1 {
		// TODO
		return fmt.Errorf("too many arguments")
	}

	defer func() {
		// reset dynamic tab completions
		s.cdCompletions = []string{}
	}()

	// if there is no target, cd to root 'dir'
	if len(args) == 0 {
		s.resetcwd()
		return nil
	}

	destination := args[0]
	if destination == ".." {
		// TODO cd .. from pod goes up to namespace,
		// cd .. from namespace goes up to contexts
		s.resetcwd()
		return nil
	}

	s.namespace = &destination

	return nil
}

func (s *session) resetcwd() {
	s.namespace = nil
}

func (s *session) getCdCompletions(input string) []string {
	return s.cdCompletions
}
