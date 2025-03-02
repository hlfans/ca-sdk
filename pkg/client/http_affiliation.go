package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hlfans/ca-sdk/pkg/entity"
	"github.com/hlfans/ca-sdk/pkg/request"
	"github.com/hlfans/ca-sdk/pkg/response"
)

const (
	endpointAffiliationList   = "%s/api/v1/affiliations%s"
	endpointAffiliationCreate = "%s/api/v1/affiliations%s"
	endpointAffiliationDelete = "%s/api/v1/affiliations/%s"
)

func (c *httpClient) AffiliationList(ctx context.Context, rootAffiliation ...string) ([]entity.Identity, []entity.Affiliation, error) {
	var reqUrl string

	if len(rootAffiliation) == 1 {
		reqUrl = fmt.Sprintf(endpointAffiliationList, c.config.Host, `/`+rootAffiliation[0])
	} else {
		reqUrl = fmt.Sprintf(endpointAffiliationList, c.config.Host, ``)
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err = c.setAuthToken(req, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to set auth token: %w", err)
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to do request: %w", err)
	}

	var affiliationResponse response.AffiliationList

	if err = c.processResponse(resp, &affiliationResponse, http.StatusOK, http.StatusCreated); err != nil {
		return nil, nil, err
	}

	return affiliationResponse.Identities, affiliationResponse.Affiliations, nil
}

func (c *httpClient) AffiliationCreate(ctx context.Context, name string, opts ...AffiliationOpt) error {
	var (
		reqUrl string
		err    error
	)
	u := url.Values{}

	for _, opt := range opts {
		if err = opt(&u); err != nil {
			return err
		}
	}

	if v := u.Encode(); v == `` {
		reqUrl = fmt.Sprintf(endpointAffiliationCreate, c.config.Host, ``)
	} else {
		reqUrl = fmt.Sprintf(endpointAffiliationCreate, c.config.Host, `?`+v)
	}

	reqBytes, err := json.Marshal(request.AddAffiliationRequest{Name: name})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if err = c.setAuthToken(req, reqBytes); err != nil {
		return fmt.Errorf("failed to set auth token: %w", err)
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	var affiliationCreateResponse response.AffiliationCreate

	if err = c.processResponse(resp, &affiliationCreateResponse, http.StatusCreated); err != nil {
		return err
	}

	return nil
}

func (c *httpClient) AffiliationDelete(ctx context.Context, name string, opts ...AffiliationOpt) ([]entity.Identity, []entity.Affiliation, error) {
	var (
		reqUrl string
		err    error
	)

	u := url.Values{}

	for _, opt := range opts {
		if err = opt(&u); err != nil {
			return nil, nil, fmt.Errorf("failed to set option: %w", err)
		}
	}

	if v := u.Encode(); v == `` {
		reqUrl = fmt.Sprintf(endpointAffiliationDelete, c.config.Host, name)
	} else {
		reqUrl = fmt.Sprintf(endpointAffiliationDelete, c.config.Host, name+`?`+v)
	}

	req, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	if err = c.setAuthToken(req, nil); err != nil {
		return nil, nil, fmt.Errorf("failed to set auth token: %w", err)
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to do request: %w", err)
	}

	var affiliationDeleteResponse response.AffiliationDelete

	if err = c.processResponse(resp, &affiliationDeleteResponse, http.StatusOK); err != nil {
		return nil, nil, err
	}
	return affiliationDeleteResponse.Identities, affiliationDeleteResponse.Affiliations, nil
}
