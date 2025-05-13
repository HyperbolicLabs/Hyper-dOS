package sh

func (s *session) cd(destination *string) {
	if s.clientset == nil {
		s.nocluster()
		return
	}

	// if destination is nil, then we are going back to the root
	if destination == nil {
		s.resetcwd()
	} else if *destination == ".." {
		// TODO cd .. from pod goes up to namespace,
		// cd .. from namespace goes up to contexts
		s.resetcwd()
	} else {
		s.namespace = destination
	}

	// reset completions
	s.cdCompletions = []string{}
}

func (s *session) resetcwd() {
	s.namespace = nil
}

func (s *session) getCdCompletions(input string) []string {
	return s.cdCompletions
}
