package sh

func (s *session) nocluster() {
	s.writeln("no cluster connected")
}
