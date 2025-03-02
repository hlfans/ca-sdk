package client

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"

	"github.com/hlfans/ca-sdk/pkg/entity"
	"github.com/hlfans/ca-sdk/pkg/request"
	"github.com/hlfans/ca-sdk/pkg/response"
)

type Client interface {
	// CAInfo Getting information about CA
	CAInfo(ctx context.Context) (*response.CAInfo, error)

	Register(ctx context.Context, req request.Registration) (string, error)
	Enroll(ctx context.Context, name, secret string, req *x509.CertificateRequest, opts ...EnrollOpt) (
		*x509.Certificate, interface{}, error)
	Revoke(ctx context.Context, req request.RevocationRequest) (*pkix.CertificateList, error)
	IdentityList(ctx context.Context) ([]entity.Identity, error)
	IdentityGet(ctx context.Context, enrollId string) (*entity.Identity, error)
	CertificateList(ctx context.Context, opts ...CertificateListOpt) ([]*x509.Certificate, error)
	// AffiliationList lists all affiliations and identities of identity affiliation
	AffiliationList(ctx context.Context, rootAffiliation ...string) ([]entity.Identity, []entity.Affiliation, error)
	AffiliationCreate(ctx context.Context, name string, opts ...AffiliationOpt) error
	AffiliationDelete(ctx context.Context, name string, opts ...AffiliationOpt) ([]entity.Identity, []entity.Affiliation, error)
}
