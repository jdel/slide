package keys

import (
	"net"

	"github.com/nats-io/nkeys"
	"golang.org/x/crypto/ssh/agent"
)

const seedLen = 32

// Just wipe slice with 'x', for clearing contents of creds or nkey seed file.
func wipeSlice(buf []byte) {
	for i := range buf {
		buf[i] = 'x'
	}
}

func GetKeyPairsFromAgent(socket string) ([]nkeys.KeyPair, error) {
	agentConn, err := net.Dial("unix", socket)
	if err != nil {
		return nil, err
	}

	agentClient := agent.NewClient(agentConn)
	signerList, err := agentClient.Signers()
	if err != nil {
		return nil, err
	}

	var keyPairs []nkeys.KeyPair
	for _, signer := range signerList {
		if signer.PublicKey().Type() == "ssh-ed25519" {
			kp, err := FromSigner(signer)
			if err != nil {
				return nil, err
			}
			keyPairs = append(keyPairs, kp)
		}
	}
	if len(keyPairs) == 0 {
		return nil, ErrNoKeyFound
	}
	return keyPairs, nil
}
