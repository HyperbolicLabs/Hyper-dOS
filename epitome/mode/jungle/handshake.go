package jungle

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type handshakeResponse struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ClusterName  string `json:"cluster_name"`
}

type Handshaker interface {
	Do(*http.Request) (*http.Response, error)
}

func (a *agent) handshake(
	handshaker Handshaker,
) (response *handshakeResponse, err error) {
	gatewayUrl := a.cfg.Default.HYPERBOLIC_GATEWAY_URL
	token := a.cfg.Default.HYPERBOLIC_TOKEN

	logrus.Infof("handshaking with gateway: %v", gatewayUrl)

	// post to gateway
	req, err := http.NewRequest(
		"POST",
		gatewayUrl.String()+"/v1/hyperweb/login",
		bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		logrus.Errorf("failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	resp, err := handshaker.Do(req)
	if err != nil {
		logrus.Errorf("failed to send handshake request: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("failed to read handshake response body: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("handshake response status: %v", resp.Status)
		return
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logrus.Infof("failed to unmarshal handshake response body: %v", string(body))
		logrus.Errorf("failed to unmarshal handshake response body: %v", err)
		return
	}

	if response.ClientID == "" || response.ClientSecret == "" || response.ClusterName == "" {
		err = errors.New("handshake response is missing required fields")
		// print the unformatted json
		logrus.Errorf("handshake response: %v", string(body))
		logrus.Errorf("handshake response header: %+v", resp.Header)
		logrus.Errorf("handshake response Status: %+v", resp.Status)
		return
	}

	logrus.Infof("handshake response: cluster name is %+v", response.ClusterName)

	return
}
