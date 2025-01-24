package hyperweb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// curl -X POST "https://api.dev-hyperbolic.xyz/v1/hyperweb/register_cluster" \
// -H "Content-Type: application/json" \
// -H "Authorization: Bearer {{.Token}}" \
// -d '{"cluster_name":"test"}'

// {"success":true}
// or
// Internal Server Error

type registerResponse struct {
	Success bool `json:"success"`
}

func register(
	gatewayUrl string,
	token string,
	clusterName string,
) (*registerResponse, error) {

	// post to gateway that we have successfully bootstrapped the cluster
	// and are ready to join the Hyperbolic Supply Network

	logrus.Infof("registering cluster with gateway: %v", gatewayUrl)

	// create payload
	payload := bytes.NewBuffer(
		[]byte(
			fmt.Sprintf(
				`{"cluster_name":"%s"}`,
				clusterName,
			),
		),
	)

	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		gatewayUrl+"/v1/hyperweb/register_cluster",
		payload,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create register_cluster request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send register_cluster request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read register_cluster response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("register_cluster response status: %v, body: %v", resp.Status, string(body))
		return nil, fmt.Errorf("cluster registration failed. register_cluster response status: %v", resp.Status)
	}

	response := &registerResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		if string(body) == "Internal Server Error" {
			logrus.Errorf("got internal server error: %v", err)
			return nil, fmt.Errorf("cluster registration failed")
		} else {
			return nil, fmt.Errorf("failed to unmarshal register_cluster response: %v", err)
		}
	}

	logrus.Infof("register_cluster response: %+v", response)

	return response, nil
}
