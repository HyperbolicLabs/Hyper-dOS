package hyperweb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func handshake(
	gatewayUrl string,
	token string,
) (response *handshakeResponse, err error) {
	// post to gateway
	client := &http.Client{}
	req, err := http.NewRequest("POST", gatewayUrl+"/v1/hyperweb/login", bytes.NewBuffer([]byte(token)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return
	}

	if response.ClientID == "" || response.ClientSecret == "" || response.ClusterName == "" {
		err = errors.New("handshake response is missing required fields")
		return
	}

	return
}

type handshakeResponse struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	ClusterName  string `json:"clusterName"`
}
