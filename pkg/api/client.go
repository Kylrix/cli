package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/nathfavour/kylrix/cli/pkg/config"
	"github.com/pkg/errors"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Config     *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL: cfg.BaseURI,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		Config: cfg,
	}
}

func (c *Client) Execute(method, path string, body interface{}, result interface{}) error {
	respData, err := c.Request(method, path, body)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respData, result); err != nil {
			return errors.Wrap(err, "failed to unmarshal response")
		}
	}

	return nil
}

func (c *Client) Request(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal request body")
		}
		bodyReader = bytes.NewBuffer(data)
	}

	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")
	if c.Config.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Config.Token))
	}
	if c.Config.APIKey != "" {
		req.Header.Set("X-Kylrix-API-Key", c.Config.APIKey)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
