package user_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/rest"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ UserService = &client{}

type client struct {
	base     rest.BaseClient
	resource string
}

func NewService(baseURL string, resource string, logger logging.Logger) *client {
	c := client{
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
		resource: resource,
	}

	return &c
}

type UserService interface {
	SignIn(ctx context.Context, login, password string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	Create(ctx context.Context, dto *CreateUserDTO) (*User, error)
	Update(ctx context.Context, id uint, dto *UpdateUserDTO) error
	Delete(ctx context.Context, id uint) error
}

func (c *client) SignIn(ctx context.Context, login, password string) (*User, error) {

	c.base.Logger.Debug("adding login and password to filter options...")
	filters := []rest.FilterOptions{
		{
			Field:  "login",
			Values: []string{login},
		},
		{
			Field:  "password",
			Values: []string{password},
		},
	}

	c.base.Logger.Debug("building url with resource and filter...")
	uri, err := c.base.BuildURL(c.resource, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL due to error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request...")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending created request...")
	// TODO implement circuit breaker pattern (i. e. hystrix lib)
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode,
			response.Error.Message,
			response.Error.DeveloperMessage)
	}

	u := &User{}
	if err = json.NewDecoder(response.Body()).Decode(&u); err != nil {
		return nil, fmt.Errorf("failed to decode body due to error: %w", err)
	}

	return u, nil
}

func (c *client) GetByID(ctx context.Context, id uint) (*User, error) {

	c.base.Logger.Debug("building url with resource and filter...")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%d", c.resource, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL due to error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request...")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending created request...")
	// TODO implement circuit breaker pattern (i. e. hystrix lib)
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode,
			response.Error.Message,
			response.Error.DeveloperMessage)
	}

	defer response.Body().Close()

	u := &User{}

	err = json.NewDecoder(response.Body()).Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body due to error: %w", err)
	}

	return u, nil
}

func (c *client) Create(ctx context.Context, dto *CreateUserDTO) (*User, error) {

	c.base.Logger.Debug("building url with resource and filter...")
	uri, err := c.base.BuildURL(c.resource, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL. error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("marshaling dto to bytes...")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dto due to err: %w", err)
	}

	c.base.Logger.Debug("creating new request...")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending created request")
	// TODO implement circuit breaker pattern (i. e. hystrix lib)
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode,
			response.Error.Message,
			response.Error.DeveloperMessage)
	}

	c.base.Logger.Debug("parsing location header...")
	userURL, err := response.Location()
	if err != nil {
		return nil, fmt.Errorf("failed to get Location header due to error: %w", err)
	}
	c.base.Logger.Tracef("Location: %s", userURL.String())

	splitCategoryURL := strings.Split(userURL.String(), "/")
	userIDStr := splitCategoryURL[len(splitCategoryURL)-1]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, err
	}

	u, err := c.GetByID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (c *client) Update(ctx context.Context, id uint, dto *UpdateUserDTO) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%d", c.resource, id), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("marshaling dto to bytes...")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal dto due to err: %w", err)
	}

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest(http.MethodPatch, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}
	if !response.IsOk {
		return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}

	return nil
}

func (c *client) Delete(ctx context.Context, id uint) error {
	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%d", c.resource, id), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("create new request")
	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("send request")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}

	return nil
}
