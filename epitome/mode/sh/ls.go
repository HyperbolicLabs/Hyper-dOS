package sh

func (s *session) ls(target *string) {
	if target != nil {
		s.writeln("TODO: ls <target> not implemented")
		return
	}

	// if we are not in a namespace, list namespaces
	if s.namespace == nil {
		s.listNamespaces(s.clientset)
	} else {
		// if we are in a namespace, list pods
		s.listPods(s.clientset, *s.namespace)
	}

}
