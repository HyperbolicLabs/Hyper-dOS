package maintain

import (
	"context"
	"time"

	argo "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"
)

func (a *agent) updateBaronConditions() error {
	// first, update the base condition of "maintainance robot is active"

	conditions := []argo.ApplicationCondition{
		{
			Type:    "RobotActive",
			Message: "Epitome Maintainance Robot is Active",
			LastTransitionTime: &metav1.Time{
				Time: time.Now(),
			},
		},
	}

	// check if there is an error with the ceph cluster
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	err := a.checkCephClusterHealth(
		ctx,
		types.NamespacedName{
			Namespace: "rook-ceph-external",
			Name:      "rook-ceph-external",
		},
	)
	if err != nil {
		conditions = append(conditions, argo.ApplicationCondition{
			Type:    "Error",
			Message: err.Error(),
			LastTransitionTime: &metav1.Time{
				Time: time.Now(),
			},
		})
	}

	ctx, cancel = context.WithTimeout(
		context.Background(), 2*time.Second)
	defer cancel()

	// set the condition on the hyperdos app in the argocd namespace
	app, err := a.argoClient.
		ArgoprojV1alpha1().
		Applications("argocd").
		Get(ctx, "hyperdos", metav1.GetOptions{})
	if err != nil {
		return err
	}

	// WARNING: this will clobber
	// so if argo ever decides to build something using
	// the conditions functionality, we will break it.
	// I would like to move to a hyperdos CRD before that happens
	app.Status.Conditions = conditions

	// note that this produces some log spam on older
	// installations of argo :(
	_, err = a.argoClient.ArgoprojV1alpha1().
		Applications("argocd").
		Update(ctx, app, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
