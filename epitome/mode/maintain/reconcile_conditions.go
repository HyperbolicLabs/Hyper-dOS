package maintain

import (
	"context"
	"time"

	argo "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) updateBaronConditions() error {
	// first, update the base condition of "maintainance robot is active"

	exampleConditions := []argo.ApplicationCondition{
		{
			Type:    "RobotActive",
			Message: "Epitome Maintainance Robot is Active",
			LastTransitionTime: &metav1.Time{
				Time: time.Now(),
			},
		},
		{
			Type:    "Error",
			Message: "ExampleError",
			LastTransitionTime: &metav1.Time{
				Time: time.Now(),
			},
		},
	}

	ctx, cancel := context.WithTimeout(
		context.Background(), 5*time.Second)
	defer cancel()

	// set the condition on the hyperdos app in the argocd namespace
	app, err := a.argoClient.
		ArgoprojV1alpha1().
		Applications("argocd").
		Get(ctx, "hyperdos", metav1.GetOptions{})
	if err != nil {
		return err
	}

	newConditions := exampleConditions

	// WARNING: this will clobber
	// so if argo ever decides to build something using
	// the conditions functionality, we will break it.
	// I would like to move to a hyperdos CRD before that happens
	app.Status.Conditions = newConditions

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
