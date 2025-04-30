package keys

import (
	"crypto/ed25519"
	"crypto/rand"
	"io"
	"os"

	"github.com/nats-io/nkeys"
	"golang.org/x/crypto/ssh"
)

type fkp struct {
	nk    nkeys.KeyPair
	curve nkeys.KeyPair
}

// FromFile creates a hybrid nkey from a private ed25519 SSH Key with a x25519 curve component used for Seal and Open
func FromFile(sshKeyFile, passphrase string) (nkeys.KeyPair, error) {
	var kp fkp

	pk, err := getRawPrivateKeyFromFile(sshKeyFile, passphrase)
	if err != nil {
		return nil, err
	}

	nk, err := nkeys.FromRawSeed(nkeys.PrefixByteUser, pk.Seed())
	if err != nil {
		return nil, err
	}
	curveSeed, err := nkeys.EncodeSeed(nkeys.PrefixByteCurve, pk.Seed())
	defer wipeSlice(curveSeed)
	if err != nil {
		return nil, err
	}
	curve, err := nkeys.FromCurveSeed(curveSeed)
	if err != nil {
		return nil, err
	}

	kp.nk = nk
	kp.curve = curve
	return &kp, nil
}

// FromCreds creates a hybrid nkey from a NATS creds file with a x25519 curve component used for Seal and Open
func FromCreds(file string) (nkeys.KeyPair, error) {
	var kp fkp

	contents, err := getRawBytesFromFile(file)
	defer wipeSlice(contents)
	if err != nil {
		return nil, err
	}

	nk, err := nkeys.ParseDecoratedUserNKey(contents)
	if err != nil {
		return nil, err
	}

	seed, err := nk.Seed()
	if err != nil {
		return nil, err
	}

	_, rawSeed, err := nkeys.DecodeSeed(seed)
	defer wipeSlice(rawSeed)
	if err != nil {
		return nil, err
	}

	pk := ed25519.NewKeyFromSeed(rawSeed)
	curveSeed, err := nkeys.EncodeSeed(nkeys.PrefixByteCurve, pk.Seed())
	if err != nil {
		return nil, err
	}
	curve, err := nkeys.FromCurveSeed(curveSeed)
	if err != nil {
		return nil, err
	}

	kp.nk = nk
	kp.curve = curve
	return &kp, nil
}

func (pair *fkp) Seed() ([]byte, error) {
	return pair.nk.Seed()
}

func (pair *fkp) PublicKey() (string, error) {
	return pair.nk.PublicKey()
}

func (pair *fkp) PrivateKey() ([]byte, error) {
	return pair.nk.PrivateKey()
}

func (pair *fkp) Sign(input []byte) ([]byte, error) {
	return pair.nk.Sign(input)
}

func (pair *fkp) Seal(input []byte, _ string) ([]byte, error) {
	// Override recipient, use the curve key public key
	recipient, err := pair.curve.PublicKey()
	if err != nil {
		return nil, err
	}
	// pass through to SealWithRand
	return pair.SealWithRand(input, recipient, rand.Reader)
}

func (pair *fkp) SealWithRand(input []byte, recipient string, rr io.Reader) ([]byte, error) {
	return pair.curve.SealWithRand(input, recipient, rand.Reader)
}

func (pair *fkp) Open(input []byte, _ string) ([]byte, error) {
	// Override sender, use the curve key public key
	sender, err := pair.curve.PublicKey()
	if err != nil {
		return nil, err
	}
	return pair.curve.Open(input, sender)
}

func (pair *fkp) Wipe() {
	pair.nk.Wipe()
	pair.curve.Wipe()
}

func (pair *fkp) Verify(input []byte, sig []byte) error {
	return pair.nk.Verify(input, sig)
}

func getRawBytesFromFile(filename string) ([]byte, error) {
	fh, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	b, err := io.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func getRawPrivateKeyFromFile(sshKey string, passphrase string) (*ed25519.PrivateKey, error) {
	pemBytes, err := getRawBytesFromFile(sshKey)
	if err != nil {
		return nil, err
	}
	defer wipeSlice(pemBytes)

	rawSigner, err := ssh.ParseRawPrivateKey(pemBytes)

	switch err.(type) {
	case *ssh.PassphraseMissingError:
		if passphrase == "" {
			return nil, ErrSshKeyRequiresPassphrase
		}

		rawSigner, err = ssh.ParseRawPrivateKeyWithPassphrase(pemBytes, []byte(passphrase))
		if err != nil {
			return nil, err
		}
	default:
		if err != nil {
			return nil, err
		}
	}

	privateKey, ok := rawSigner.(*ed25519.PrivateKey)
	if !ok {
		return nil, ErrWrongKeyType
	}
	return privateKey, nil
}
