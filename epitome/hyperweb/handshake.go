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
) (clientID *string, clientSecret *string, err error) {
	// -H "Content-Type: application/json" \
	// -H "Authorization: Bearer $token" \

	// post to gateway
	client := &http.Client{}
	req, err := http.NewRequest("POST", gatewayUrl+"/v1/hyperweb/login", bytes.NewBuffer([]byte(token)))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if err = json.Unmarshal(body, &handshakeResponse); err != nil {
		return nil, nil, err
	}

	if handshakeResponse.ClientID == "" || handshakeResponse.ClientSecret == "" {
		return nil, nil, errors.New("handshake response is missing required fields")
	}

	return &handshakeResponse.ClientID, &handshakeResponse.ClientSecret, nil
}

var handshakeResponse struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
}
