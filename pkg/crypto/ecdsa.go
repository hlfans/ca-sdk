package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"reflect"
)

var (
	// precomputed curves half order values for efficiency
	ecCurveHalfOrders = map[elliptic.Curve]*big.Int{
		elliptic.P224(): new(big.Int).Rsh(elliptic.P224().Params().N, 1),
		elliptic.P256(): new(big.Int).Rsh(elliptic.P256().Params().N, 1),
		elliptic.P384(): new(big.Int).Rsh(elliptic.P384().Params().N, 1),
		elliptic.P521(): new(big.Int).Rsh(elliptic.P521().Params().N, 1),
	}
)

type ecdsaSigner struct {
	key  *ecdsa.PrivateKey
	cert *x509.Certificate
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (e *ecdsaSigner) Public() crypto.PublicKey {
	return e.key.Public()
}

func (e *ecdsaSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	R, S, err := ecdsa.Sign(rand, e.key, digest)
	if err != nil {
		return nil, fmt.Errorf("sign message: %w", err)
	} else {
		preventMalleability(e.key, S)
	}

	signature, err = asn1.Marshal(ecdsaSignature{R, S})
	if err != nil {
		return nil, fmt.Errorf("marshal asn1 signature: %w", err)
	}
	return
}

func (e *ecdsaSigner) Certificate() []byte {
	b := new(bytes.Buffer)
	if err := pem.Encode(b, &pem.Block{Type: "CERTIFICATE", Bytes: e.cert.Raw}); err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewSigner(cert *x509.Certificate, key interface{}) (Signer, error) {
	ecdsaKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf(`invalid key type; expected ecdsa, got %s`, reflect.TypeOf(key))
	}
	return &ecdsaSigner{key: ecdsaKey, cert: cert}, nil
}

func NewPrivateKey() (crypto.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// from gohfc
func preventMalleability(k *ecdsa.PrivateKey, S *big.Int) {
	halfOrder := ecCurveHalfOrders[k.Curve]
	if S.Cmp(halfOrder) == 1 {
		S.Sub(k.Params().N, S)
	}
}
