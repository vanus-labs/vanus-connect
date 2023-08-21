package vanusai

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient *resty.Client
	url        string
}

func NewClient(url string) *Client {
	return &Client{
		httpClient: resty.New(),
		url:        url,
	}
}

func (c *Client) Chat(ctx context.Context, req *ChatRequest) (string, error) {
	return c.ChatWithUrl(ctx, c.url, req)
}

func (c *Client) ChatWithUrl(ctx context.Context, urlStr string, req *ChatRequest) (string, error) {
	resp, err := c.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "text/plain").
		SetHeader("X-Vanusai-Sessionid", req.SessionID).
		SetHeader("X-Vanusai-Model", req.Model).
		SetBody(req).
		SetContext(ctx).
		Post(urlStr)
	if err != nil {
		return "", err
	}
	if !resp.IsSuccess() {
		return "", fmt.Errorf("ai response code:%d,body:%s", resp.StatusCode(), resp.String())
	}
	return string(resp.Body()), nil
}
