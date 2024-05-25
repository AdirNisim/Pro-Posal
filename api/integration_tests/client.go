package integrationtests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pro-posal/webserver/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ApiClient struct {
	client    *http.Client
	baseURL   string
	authToken string
}

func NewApiClient(baseURL string, email string, password string) (*ApiClient, error) {
	client := &ApiClient{
		client:  &http.Client{},
		baseURL: baseURL,
	}
	err := client.ObtainAuthToken(email, password)
	return client, err
}

func (c *ApiClient) ObtainAuthToken(email string, password string) error {
	requestBody := map[string]string{
		"email":    email,
		"password": password,
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	reqPayload := bytes.NewBuffer(bodyBytes)
	req, err := http.NewRequest("POST", c.baseURL+"/users/login", reqPayload)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to obtain auth token, got %d status code", resp.StatusCode)
	}

	var authTokenResp api.PostUsersLoginResponseBody
	err = json.NewDecoder(resp.Body).Decode(&authTokenResp)
	if err != nil {
		return err
	}

	c.authToken = authTokenResp.AccessToken
	return nil
}

func (c *ApiClient) Post(t *testing.T, url string, payload any, statusCode int, resp any) {
	t.Helper()

	bodyBytes, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", c.baseURL+url, bytes.NewBuffer(bodyBytes))
	require.NoError(t, err)
	req.Header.Add("Authorization", "Bearer "+c.authToken)

	httpResp, err := c.client.Do(req)
	require.NoError(t, err)
	defer httpResp.Body.Close()

	assert.Equal(t, statusCode, httpResp.StatusCode)

	if resp != nil {
		err = json.NewDecoder(httpResp.Body).Decode(resp)
		require.NoError(t, err)
	}
}

func (c *ApiClient) Get(t *testing.T, url string, statusCode int, resp any) {
	// TODO: Implement me
}
