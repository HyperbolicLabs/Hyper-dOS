package maintain

import (
	"context"
	"fmt"
	"time"
)

func (a *agent) reconcile() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// sometimes, calico seems to behave strangely
	// and crashes containers across the whole cluster on creation/deletion.
	// a restart usually fixes things though.
	// So, in the interest of simple stability,
	// when a baron has calico installed,
	// we'll just periodically restart the daemonset (until the calico team fixes upstream)
	err := a.restartCalicoIfExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to restart calico: %v", err)
	}

	a.logger.Debug("reconcile complete")
	return nil
}
