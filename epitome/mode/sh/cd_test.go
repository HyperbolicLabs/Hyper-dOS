package sh

import (
	"testing"

	"epitome.hyperbolic.xyz/config"
	shtest "epitome.hyperbolic.xyz/mode/sh/test"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_session_cd(t *testing.T) {
	// create a new session
	s, close, err := NewSession(
		&config.Config{},
		zap.NewNop(),
		shtest.InitFakeCluster(),
	)
	if err != nil {
		t.Errorf("failed to create session: %v", err)
	}

	defer close()

	if err := s.cd("hyperdos"); err != nil {
		t.Errorf("failed to cd: %v", err)
	}

	require.Equal(t, *s.namespace, "hyperdos")
}
