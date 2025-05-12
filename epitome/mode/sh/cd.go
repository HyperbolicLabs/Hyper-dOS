package sh

func (s *session) cd(destination *string) {
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
}

func (s *session) resetcwd() {
	s.namespace = nil
}
