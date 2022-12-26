package lot_service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/levelord1311/backendForSharedProject/api_service/internal/apperror"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/logging"
	"github.com/levelord1311/backendForSharedProject/api_service/pkg/rest"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ LotService = &client{}

type client struct {
	Resource string
	base     rest.BaseClient
}

func NewService(baseURL string, resource string, logger logging.Logger) *client {
	return &client{
		Resource: resource,
		base: rest.BaseClient{
			BaseURL: baseURL,
			HTTPClient: &http.Client{
				Timeout: 10 * time.Second,
			},
			Logger: logger,
		},
	}
}

type LotService interface {
	GetByUserID(ctx context.Context, id string) ([]byte, error)
	GetByLotID(ctx context.Context, id string) ([]byte, error)
	GetWithFilter(ctx context.Context, rQuery string) ([]byte, error)
	Create(ctx context.Context, dto *CreateLotDTO) (uint, error)
	Update(ctx context.Context, dto *UpdateLotDTO) error
	Delete(ctx context.Context, lotID, userID string) error
}

func (c *client) GetByUserID(ctx context.Context, id string) ([]byte, error) {

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s/%s", c.Resource, "/user", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending request..")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}

	c.base.Logger.Debug("reading response body..")
	lots, err := response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return lots, nil
}

func (c *client) GetByLotID(ctx context.Context, id string) ([]byte, error) {

	c.base.Logger.Debug("build url with resource and filter")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s/%s", c.Resource, "/lot", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL. error: %v", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending request..")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}

	c.base.Logger.Debug("reading response body..")
	lot, err := response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}
	return lot, nil
}

func (c *client) GetWithFilter(ctx context.Context, rQuery string) ([]byte, error) {

	c.base.Logger.Debug("build url with resource and raw query")
	c.base.Logger.Debug("c.RESOURCE:", c.Resource)
	c.base.Logger.Debug("rQuery:", rQuery)

	//uri, err := c.base.BuildURL(fmt.Sprintf("%s?%s", c.Resource, rQuery), nil)
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL. error: %v", err)
	}
	uri = fmt.Sprintf("%s?%s", uri, rQuery)
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending request..")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return nil, apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}

	c.base.Logger.Debug("reading response body..")
	lots, err := response.ReadBody()
	if err != nil {
		return nil, fmt.Errorf("failed to read body")
	}

	return lots, nil
}

func (c *client) Create(ctx context.Context, dto *CreateLotDTO) (uint, error) {
	c.base.Logger.Debug("building url with resource and filter..")
	uri, err := c.base.BuildURL(c.Resource, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to build URL. error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("marshaling dto to bytes..")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal dto due to err: %w", err)
	}

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return 0, fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending created request..")
	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)
	response, err := c.base.SendRequest(req)
	if err != nil {
		return 0, fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return 0, apperror.APIError(response.Error.ErrorCode,
			response.Error.Message,
			response.Error.DeveloperMessage)
	}
	c.base.Logger.Debug("parsing location header..")
	lotURL, err := response.Location()
	if err != nil {
		return 0, fmt.Errorf("failed to get Location header")
	}
	c.base.Logger.Tracef("Location: %s", lotURL.String())

	splitCategoryURL := strings.Split(lotURL.String(), "/")
	lotIDStr := splitCategoryURL[len(splitCategoryURL)-1]
	lotID, err := strconv.Atoi(lotIDStr)
	if err != nil {
		return 0, err
	}

	return uint(lotID), nil
}

func (c *client) Update(ctx context.Context, dto *UpdateLotDTO) error {
	c.base.Logger.Debug("building url with resource and filter..")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%d", c.Resource, dto.ID), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("marshaling dto to bytes..")
	dataBytes, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal dto due to err: %w", err)
	}

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodPatch, uri, bytes.NewBuffer(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}

	c.base.Logger.Debug("sending created request..")
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

func (c *client) Delete(ctx context.Context, lotID, userID string) error {
	c.base.Logger.Debug("building url with resource and filter..")
	uri, err := c.base.BuildURL(fmt.Sprintf("%s/%s", c.Resource, lotID), nil)
	if err != nil {
		return fmt.Errorf("failed to build URL. error: %w", err)
	}
	c.base.Logger.Tracef("url: %s", uri)

	c.base.Logger.Debug("creating new request..")
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("failed to create new request due to error: %w", err)
	}
	req.Header.Set("user_id", userID)

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req = req.WithContext(reqCtx)

	c.base.Logger.Debug("sending created request..")
	response, err := c.base.SendRequest(req)
	if err != nil {
		return fmt.Errorf("failed to send request due to error: %w", err)
	}

	if !response.IsOk {
		return apperror.APIError(response.Error.ErrorCode, response.Error.Message, response.Error.DeveloperMessage)
	}
	return nil
}
