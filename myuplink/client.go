package myuplink

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/webdevops/go-common/log/slogger"
	resty "resty.dev/v3"
)

type (
	Client struct {
		logger *slogger.Logger
		http   *resty.Client

		clientId     string
		clientSecret string

		token *AuthToken
	}
)

func NewClient(logger *slogger.Logger) *Client {
	c := Client{
		logger: logger,
		http:   resty.New(),
	}

	c.http.SetRetryCount(5)
	c.http.SetRetryWaitTime(5 * time.Second)
	c.http.SetRetryMaxWaitTime(30 * time.Second)
	c.http.EnableRetryDefaultConditions()

	return &c
}

func (c *Client) SetDebugMode(val bool) {
	c.http.SetDebug(val)
}

func (c *Client) SetApiUrl(val string) {
	c.http.SetBaseURL(val)
}

func (c *Client) SetUserAgent(val string) {
	c.http.SetHeader("User-Agent", val)
}

func (c *Client) SetAuth(clientId, clientSecret string) {
	c.clientId = clientId
	c.clientSecret = clientSecret
}

func (c *Client) Connect(ctx context.Context) error {
	c.logger.Debug(`refresh myUplink auth token`)

	payload := map[string]string{
		"grant_type": "client_credentials",
		"scope":      "READSYSTEM",
	}

	result := &AuthToken{}

	resp, err := c.http.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(payload).
		SetBasicAuth(c.clientId, c.clientSecret).
		SetResult(result).
		Post("oauth/token")
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return fmt.Errorf(`unable to login to myUplink, expected status code 200, got %v`, resp.Status())
	}

	c.token = result.Process()

	return nil
}

func (c *Client) createRequest(ctx context.Context) *resty.Request {
	if c.token.IsExpired() {
		if err := c.Connect(ctx); err != nil {
			c.logger.Fatal(`unable to refresh auth token`, slog.Any("error", err))
		}

	}

	return c.http.R().SetContext(ctx).SetAuthToken(c.token.AccessToken)
}

func (c *Client) GetSystems(ctx context.Context) (*ResultSystems, error) {
	result := &ResultSystems{}

	c.logger.Debug(`fetch myUplink system list`)

	resp, err := c.createRequest(ctx).SetResult(result).Get("/v2/systems/me?page=1&itemsPerPage=100")
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf(`unable to fetch system list from myUplink, expected status code 200, got %v`, resp.Status())
	}

	return result, nil
}

func (c *Client) GetSystemDevicePoints(ctx context.Context, deviceId string) (*SystemDevicePoints, error) {
	result := &SystemDevicePoints{}

	c.logger.Debug(`fetch myUplink device points`, slog.String("deviceID", deviceId))

	resp, err := c.createRequest(ctx).SetResult(result).Get(fmt.Sprintf("/v2/devices/%s/points", deviceId))
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf(`unable to fetch device points from myUplink, expected status code 200, got %v`, resp.Status())
	}

	return result, nil
}

func (c *Client) GetSystemDeviceFirmware(ctx context.Context, deviceId string) (*SystemDeviceFirmwareInfo, error) {
	result := &SystemDeviceFirmwareInfo{}

	c.logger.Debug(`fetch myUplink device firmware info`, slog.String("deviceID", deviceId))

	resp, err := c.createRequest(ctx).SetResult(result).Get(fmt.Sprintf("/v2/devices/%s/firmware-info", deviceId))
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, fmt.Errorf(`unable to fetch device firmware info from myUplink, expected status code 200, got %v`, resp.Status())
	}

	return result, nil
}
