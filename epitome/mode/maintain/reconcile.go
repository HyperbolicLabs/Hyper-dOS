package maintain

import "github.com/sirupsen/logrus"

func (a *agent) reconcile() error {
	// patch nvidia cluster policy if we are on a buffalo baron
	// this is a hack to work around a common bug in the NVIDIA operator
	// (as only the buffalo are expected to have the NVIDIA operator installed)
	if a.cfg.Role.Buffalo {
		err := a.patchClusterPolicy()
		if err != nil {
			logrus.Errorf("failed to patch cluster policy: %v", err)
		}
	}

	// update the conditions object on the hyperdos top-level gitapp
	// these are picked up by the barons operator in the king cluster
	// and used to aggregate and display more nuanced details about the internal state
	// of the baron to suppliers
	err := a.updateBaronConditions()
	if err != nil {
		logrus.Errorf("failed to update baron conditions: %v", err)
	}

	a.logger.Debug("reconcile complete")
	return nil
}
