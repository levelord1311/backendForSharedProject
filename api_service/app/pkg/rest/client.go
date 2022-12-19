package rest

import (
	"encoding/json"
	"fmt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"net/http"
	"net/url"
	"path"
)

type BaseClient struct {
	BaseURL    string
	HTTPClient *http.Client
	Logger     logging.Logger
}

func (c *BaseClient) SendRequest(r *http.Request) (*APIResponse, error) {
	if c.HTTPClient == nil {
		return nil, ErrNoHTTPClient
	}

	r.Header.Set("Accept", "application/json; charset=utf-8")
	r.Header.Set("Content-Type", "application/json; charset=utf-8")

	response, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to send request. error: %w", err)
	}

	apiResponse := APIResponse{
		IsOk:     true,
		response: response,
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		apiResponse.IsOk = false

		defer response.Body.Close()

		var apiErr APIError
		if err = json.NewDecoder(response.Body).Decode(&apiErr); err == nil {
			apiResponse.Error = apiErr
		}
	}

	return &apiResponse, nil
}

// BuildURL expects for subResource exactly one argument or none. Other than first will be ignored.
func (c *BaseClient) BuildURL(resource string, filters []FilterOptions, subResource ...string) (string, error) {
	var resultURL string
	parsedURL, err := url.ParseRequestURI(c.BaseURL)
	if err != nil {
		return resultURL, fmt.Errorf("failed to parse base URL. error: %w", err)
	}
	parsedURL.Path = path.Join(parsedURL.Path, resource)
	if len(subResource) > 0 {
		parsedURL.Path = path.Join(parsedURL.Path, subResource[0])
	}

	if len(filters) > 0 {
		q := parsedURL.Query()
		for _, fo := range filters {
			q.Set(fo.Field, fo.ToStringWF())
		}
		parsedURL.RawQuery = q.Encode()
	}

	return parsedURL.String(), nil
}

func (c *BaseClient) Close() error {
	c.HTTPClient = nil
	return nil
}
