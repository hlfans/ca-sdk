package test

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"testing"

	"github.com/hlfans/ca-sdk/pkg/client"
	"github.com/hlfans/ca-sdk/pkg/config"
	"github.com/hlfans/ca-sdk/pkg/crypto"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type FabricCACommonSuite struct {
	suite.Suite
}

func (s *FabricCACommonSuite) TestIteration(t provider.T) {
	var (
		adminCert *x509.Certificate
		adminKey  *ecdsa.PrivateKey
	)
	t.WithNewStep("Enroll CA certificate with non admin client instance", func(sCtx provider.StepCtx) {
		var cli client.Client
		sCtx.WithNewStep("Create client instance", func(sCtx provider.StepCtx) {
			var err error
			cli, err = client.NewHttp(client.WithRawConfig(&config.CAConfig{
				Host: "http://localhost:7054",
				Tls:  config.TlsConfig{},
			}))
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(cli)
		})
		sCtx.WithNewStep("Get CA info", func(sCtx provider.StepCtx) {
			info, err := cli.CAInfo(context.Background())
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(info)
		})

		sCtx.WithNewStep("Enroll CA admin certificate", func(sCtx provider.StepCtx) {
			var (
				err error
				key interface{}

				ok bool
			)
			ctx := context.Background()
			adminCert, key, err = cli.Enroll(ctx, "admin", "adminpw", &x509.CertificateRequest{
				Subject: pkix.Name{
					CommonName: "admin",
				},
			})
			sCtx.Require().NoError(err)
			sCtx.Require().NotNil(adminCert)

			adminKey, ok = key.(*ecdsa.PrivateKey)
			sCtx.Require().True(ok)
		})
	})
	t.WithNewStep("List registered identities", func(sCtx provider.StepCtx) {
		var cli client.Client
		adminSigner, err := crypto.NewSigner(adminCert, adminKey)
		sCtx.Require().NoError(err)

		sCtx.WithNewStep("Create admin-based client instance", func(sCtx provider.StepCtx) {
			cli, err = client.NewHttp(client.WithRawConfig(&config.CAConfig{
				Host: "http://localhost:7054",
			}), client.WithIdentity(adminSigner))
			sCtx.Require().NoError(err)
		})

		sCtx.WithNewStep("List registered identities", func(sCtx provider.StepCtx) {
			certs, err := cli.CertificateList(context.Background())
			sCtx.Require().NoError(err)
			sCtx.Require().Greater(len(certs), 0)
		})
	})
}

func TestFabricCA(t *testing.T) {
	suite.RunSuite(t, new(FabricCACommonSuite))
}
