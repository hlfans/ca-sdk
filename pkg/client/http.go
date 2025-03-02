package client

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hlfans/ca-sdk/pkg/config"
	"github.com/hlfans/ca-sdk/pkg/crypto"
	"github.com/hlfans/ca-sdk/pkg/entity"
	"github.com/hlfans/ca-sdk/pkg/request"
	"github.com/hlfans/ca-sdk/pkg/response"
	"gopkg.in/yaml.v3"
)

type HttpOpt func(c *httpClient) error

// WithYamlConfig allows using YAML config from file
func WithYamlConfig(path string) HttpOpt {
	return func(c *httpClient) error {
		if configBytes, err := os.ReadFile(path); err != nil {
			return fmt.Errorf(`open config: %w`, err)
		} else {
			c.config = new(config.CAConfig)
			if err = yaml.Unmarshal(configBytes, c.config); err != nil {
				return fmt.Errorf(`unmarshal config: %w`, err)
			}
		}
		return nil
	}
}

func WithBytesConfig(configBytes []byte) HttpOpt {
	return func(c *httpClient) error {
		if err := yaml.Unmarshal(configBytes, c.config); err != nil {
			return fmt.Errorf(`unmarshal YAML config: %w`, err)
		}
		return nil
	}
}

func WithRawConfig(conf *config.CAConfig) HttpOpt {
	return func(c *httpClient) error {
		c.config = conf
		return nil
	}
}

func WithHTTPClient(client *http.Client) HttpOpt {
	return func(c *httpClient) error {
		c.client = client
		return nil
	}
}

func WithIdentity(signer crypto.Signer) HttpOpt {
	return func(c *httpClient) error {
		c.signer = signer
		return nil
	}
}

type httpClient struct {
	config *config.CAConfig
	client *http.Client
	signer crypto.Signer
}

func (c *httpClient) Register(ctx context.Context, req request.Registration) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c *httpClient) Revoke(ctx context.Context, req request.RevocationRequest) (*pkix.CertificateList, error) {
	//TODO implement me
	panic("implement me")
}

func (c *httpClient) IdentityList(ctx context.Context) ([]entity.Identity, error) {
	//TODO implement me
	panic("implement me")
}

func (c *httpClient) IdentityGet(ctx context.Context, enrollId string) (*entity.Identity, error) {
	//TODO implement me
	panic("implement me")
}

func NewHttp(opts ...HttpOpt) (Client, error) {
	var err error

	var cli httpClient

	for _, opt := range opts {
		if err = opt(&cli); err != nil {
			return nil, fmt.Errorf(`apply ca.Client option: %w`, err)
		}
	}

	if cli.config == nil {
		return nil, fmt.Errorf(`config is empty`)
	}

	if cli.client == nil {
		cli.client = http.DefaultClient
	}

	return &cli, nil
}

func (c *httpClient) createAuthToken(method string, url string, request []byte) (string, error) {
	bodyEncoded := base64.StdEncoding.EncodeToString(request)
	certEncoded := base64.StdEncoding.EncodeToString(c.signer.Certificate())
	urlEncoded := base64.URLEncoding.EncodeToString([]byte(url))

	payload := strings.Join([]string{method, urlEncoded, bodyEncoded, certEncoded}, ".")

	hasher := sha256.New()
	hasher.Write([]byte(payload))

	signature, err := c.signer.Sign(rand.Reader, hasher.Sum(nil), nil)
	if err != nil {
		return "", fmt.Errorf(`sign payload: %w`, err)
	}
	return strings.Join([]string{certEncoded, base64.StdEncoding.EncodeToString(signature)}, "."), nil
}

func (c *httpClient) setAuthToken(req *http.Request, body []byte) error {
	if token, err := c.createAuthToken(req.Method, req.URL.Path, body); err != nil {
		return fmt.Errorf("failed to create auth token: %w", err)
	} else {
		req.Header.Add(`Authorization`, token)
	}
	return nil
}

func (c *httpClient) processResponse(resp *http.Response, out interface{}, expectedHTTPStatuses ...int) error {
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if !c.expectedHTTPStatus(resp.StatusCode, expectedHTTPStatuses...) {
		return ErrUnexpectedHTTPStatus{Status: resp.StatusCode, Body: body}
	}

	var caResp response.Response
	if err = json.Unmarshal(body, &caResp); err != nil {
		return fmt.Errorf("unmarshal JSON response: %w", err)
	}

	if !caResp.Success {
		return ResponseError{Errors: caResp.Errors}
	}

	if err = json.Unmarshal(caResp.Result, out); err != nil {
		return fmt.Errorf("unmarshal result: %w", err)
	}

	return nil
}

func (c *httpClient) expectedHTTPStatus(status int, expected ...int) bool {
	for _, s := range expected {
		if s == status {
			return true
		}
	}
	return false
}
