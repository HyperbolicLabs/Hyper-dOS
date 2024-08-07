package hyperweb

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

func handshake(
	gatewayUrl string,
	token string,
) (response *handshakeResponse, err error) {
	logrus.Infof("handshaking with gateway: %v", gatewayUrl)

	// post to gateway
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		gatewayUrl+"/v1/hyperweb/login",
		// gatewayUrl,
		bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		logrus.Errorf("failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	resp, err := client.Do(req)
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

	logrus.Infof("handshake response: %+v", response)

	return
}
