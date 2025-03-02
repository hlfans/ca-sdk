package crypto

import "crypto"

type Signer interface {
	crypto.Signer
	Certificate() []byte
}
