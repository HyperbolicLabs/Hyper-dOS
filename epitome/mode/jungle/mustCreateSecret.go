package jungle

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) createTailscaleOperatorOAuthSecret(
	clientId string,
	clientSecret string,
) error {
	// create secret
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "operator-oauth", // https://tailscale.com/kb/1306/gitops-acls-github?q=operator-oauth%20secret
			Namespace: a.cfg.HyperwebNamespace,
		},
		Data: map[string][]byte{
			"client_id":     []byte(clientId),
			"client_secret": []byte(clientSecret),
		},
	}

	_, err := a.clientset.CoreV1().Secrets(a.cfg.HyperwebNamespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create operator-oauth secret: %v", err)
	}

	logrus.Infof("created operator-oauth secret")

	// TODO delete hyperbolic token?
	return nil
}
