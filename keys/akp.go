package keys

import (
	"crypto/rand"
	"io"

	"github.com/nats-io/nkeys"
	"golang.org/x/crypto/ssh"
)

type akp struct {
	signer ssh.Signer
	nk     nkeys.KeyPair
}

// FromSigner creates a hybrid nkey from a public ed25519 SSH Key with a x25519 curve component used for Seal and Open
func FromSigner(signer ssh.Signer) (nkeys.KeyPair, error) {
	var kp akp
	kp.signer = signer

	rawPublicKey := signer.PublicKey().Marshal()
	seed := rawPublicKey[len(rawPublicKey)-seedLen:]
	key, err := nkeys.Encode(nkeys.PrefixByteUser, seed)
	if err != nil {
		return nil, err
	}

	nk, err := nkeys.FromPublicKey(string(key))
	if err != nil {
		return nil, err
	}
	kp.nk = nk

	return &kp, nil
}

func (pair *akp) Seed() ([]byte, error) {
	return pair.nk.Seed()
}

func (pair *akp) PublicKey() (string, error) {
	return pair.nk.PublicKey()
}

func (pair *akp) PrivateKey() ([]byte, error) {
	// TODO: probably shouldn't return a nil priv key here
	// return nil, ErrInvalidAgentKeyOperation
	// pk, err := pair.nk.PrivateKey()
	// if err != nil {
	// 	return nil, nil
	// }
	// return pk, nil
	return pair.nk.PrivateKey()
}

func (pair *akp) Sign(input []byte) ([]byte, error) {
	signature, err := pair.signer.Sign(rand.Reader, input)
	if err != nil {
		return nil, err
	}
	return signature.Blob, nil
}

func (pair *akp) Seal(_ []byte, _ string) ([]byte, error) {
	return nil, ErrInvalidAgentKeyOperation
}

func (pair *akp) SealWithRand(_ []byte, _ string, _ io.Reader) ([]byte, error) {
	return nil, ErrInvalidAgentKeyOperation
}

func (pair *akp) Open(_ []byte, _ string) ([]byte, error) {
	return nil, ErrInvalidAgentKeyOperation
}

func (pair *akp) Wipe() {
	pair.nk.Wipe()
}

func (pair *akp) Verify(input []byte, sig []byte) error {
	return pair.nk.Verify(input, sig)
}
