package sh

import (
	"testing"

	"epitome.hyperbolic.xyz/config"
	shtest "epitome.hyperbolic.xyz/mode/sh/test"
	"go.uber.org/zap"
)

func Test_session_ls(t *testing.T) {
	s, _, err := NewSession(
		&config.Config{},
		zap.NewNop(),
		shtest.InitFakeCluster(),
	)
	if err != nil {
		t.Errorf("failed to create session: %v", err)
	}

	if err := s.ls(); err != nil {
		t.Errorf("failed to ls: %v", err)
	}

	// check that cdCompletions populated
	if len(s.cdCompletions) == 0 {
		t.Errorf("cd completions not populated")
	}

	// make sure that 'hyperdos' namespace is in cdCompletions
	found := false
	for _, c := range s.cdCompletions {
		if c == "hyperdos" {
			found = true
		}
	}
	if !found {
		t.Errorf("hyperdos namespace not in cd completions")
	}
}
