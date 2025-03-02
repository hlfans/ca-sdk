package response

import (
	"encoding/json"

	"github.com/hlfans/ca-sdk/pkg/entity"
)

type (
	Response struct {
		Success  bool            `json:"success"`
		Result   json.RawMessage `json:"result"`
		Errors   []Message       `json:"errors"`
		Messages []Message       `json:"messages"`
	}

	Message struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	CAInfo struct {
		CAName  string `json:"CAName"`
		CAChain string `json:"CAChain"`
		Version string `json:"Version"`
	}

	Registration struct {
		Secret string `json:"secret"`
	}

	Enrollment struct {
		Cert       string `json:"Cert"`
		ServerInfo CAInfo `json:"ServerInfo"`
	}

	IdentityList struct {
		Identities []entity.Identity `json:"identities"`
	}

	CertificateList struct {
		CAName string               `json:"caname"`
		Certs  []CertificateListPEM `json:"certs"`
	}

	CertificateListPEM struct {
		PEM string `json:"PEM"`
	}

	Revoke struct {
		RevokedCerts []entity.RevokedCert
		CRL          []byte
	}

	AffiliationList struct {
		Name         string               `json:"name"`
		Affiliations []entity.Affiliation `json:"affiliations"`
		Identities   []entity.Identity    `json:"identities"`
		CAName       string               `json:"caname"`
	}

	AffiliationCreate struct {
		Name   string `json:"name"`
		CAName string `json:"caname"`
	}

	AffiliationDelete struct {
		AffiliationList
	}
)
