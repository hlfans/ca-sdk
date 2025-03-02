package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hlfans/ca-sdk/pkg/response"
)

func (c *httpClient) CAInfo(ctx context.Context) (*response.CAInfo, error) {
	req, err := http.NewRequest(http.MethodGet, c.config.Host+`/api/v1/cainfo`, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("process http request: %w", err)
	}

	var caInfoResp response.CAInfo
	if err = c.processResponse(resp, &caInfoResp, http.StatusOK); err != nil {
		return nil, err
	}

	return &caInfoResp, nil
}
