package jungle

import (
	"context"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *agent) mustCreateTailscaleOperatorOAuthSecret(
	clientId string,
	clientSecret string,
) (err error) {

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

	_, err = a.clientset.CoreV1().Secrets(a.cfg.HyperwebNamespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		logrus.Fatalf("failed to create secret: %v", err)
	}

	logrus.Infof("created operator-oauth secret")

	// TODO delete hyperbolic token?
	return
}
