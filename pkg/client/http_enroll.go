package client

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"

	"github.com/cloudflare/cfssl/signer"
	"github.com/hlfans/ca-sdk/pkg/crypto"
	"github.com/hlfans/ca-sdk/pkg/response"
)

const enrollEndpoint = `/api/v1/enroll`

func (c *httpClient) Enroll(ctx context.Context, name, secret string, req *x509.CertificateRequest, opts ...EnrollOpt) (*x509.Certificate, interface{}, error) {
	var err error

	options := &EnrollOpts{}
	for _, opt := range opts {
		if err = opt(options); err != nil {
			return nil, nil, fmt.Errorf("enroll option error: %w", err)
		}
	}

	if options.PrivateKey == nil {
		if options.PrivateKey, err = crypto.NewPrivateKey(); err != nil {
			return nil, nil, fmt.Errorf(`failed to generate private key: %w`, err)
		}
	}

	if options.Profile == "" {
		options.Profile = EnrollProfileDefault
	}

	// Add default signature algorithm if not defined
	if req.SignatureAlgorithm == x509.UnknownSignatureAlgorithm {
		req.SignatureAlgorithm = x509.ECDSAWithSHA256
	}

	csr, err := x509.CreateCertificateRequest(rand.Reader, req, options.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to create CSR: %w`, err)
	}

	pemCsr := pem.EncodeToMemory(&pem.Block{Type: `CERTIFICATE REQUEST`, Bytes: csr})

	reqBytes, err := json.Marshal(signer.SignRequest{Request: string(pemCsr), Profile: string(options.Profile)})
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to marshal request: %w`, err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.config.Host+enrollEndpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to create http request: %w`, err)
	}
	httpReq.SetBasicAuth(name, secret)

	resp, err := c.client.Do(httpReq.WithContext(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf(`http request failed: %w`, err)
	}

	var enrollResp response.Enrollment

	if err = c.processResponse(resp, &enrollResp, http.StatusCreated); err != nil {
		return nil, nil, fmt.Errorf(`process response failed: %w`, err)
	}

	certDecoded, err := base64.StdEncoding.DecodeString(enrollResp.Cert)
	if err != nil {
		return nil, nil, fmt.Errorf(`decode certificate failed: %w`, err)
	}

	certBlock, _ := pem.Decode(certDecoded)
	if certBlock == nil {
		return nil, nil, fmt.Errorf(`failed to decode certificate`)
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, fmt.Errorf(`failed to parse certificate: %w`, err)
	}

	return cert, options.PrivateKey, nil
}
