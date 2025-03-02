package client

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hlfans/ca-sdk/pkg/response"
)

const endpointCertificateList = "%s/api/v1/certificates%s"

func (c *httpClient) CertificateList(ctx context.Context, opts ...CertificateListOpt) ([]*x509.Certificate, error) {
	var (
		reqUrl string
		err    error
	)

	u := url.Values{}
	for _, opt := range opts {
		if err = opt(&u); err != nil {
			return nil, fmt.Errorf("apply opt: %w", err)
		}
	}

	if v := u.Encode(); v == `` {
		reqUrl = fmt.Sprintf(endpointCertificateList, c.config.Host, ``)
	} else {
		reqUrl = fmt.Sprintf(endpointCertificateList, c.config.Host, `?`+v)
	}

	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	if err = c.setAuthToken(req, nil); err != nil {
		return nil, fmt.Errorf("set authorization token: %w", err)
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("process request: %w", err)
	}

	var certResponse response.CertificateList

	if err = c.processResponse(resp, &certResponse, http.StatusOK); err != nil {
		return nil, fmt.Errorf("process response: %w", err)
	}

	certs := make([]*x509.Certificate, len(certResponse.Certs))
	for i, v := range certResponse.Certs {
		b, _ := pem.Decode([]byte(v.PEM))
		if b == nil {
			return nil, fmt.Errorf("failed to parse PEM block: %s", v)
		}
		if cert, err := x509.ParseCertificate(b.Bytes); err != nil {
			return nil, fmt.Errorf("parse certificate: %w", err)
		} else {
			certs[i] = cert
		}
	}

	return certs, nil
}
