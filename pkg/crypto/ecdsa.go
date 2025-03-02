package crypto

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"reflect"
)

type ecdsaSigner struct {
	key  *ecdsa.PrivateKey
	cert *x509.Certificate
}

func (e *ecdsaSigner) Public() crypto.PublicKey {
	return e.key.Public()
}

func (e *ecdsaSigner) Sign(rand io.Reader, digest []byte, opts crypto.SignerOpts) (signature []byte, err error) {
	return e.key.Sign(rand, digest, opts)
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
