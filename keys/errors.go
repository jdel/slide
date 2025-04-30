package keys

const ErrInvalidAgentKeyOperation = keysError("keys: agent key is only valid for signing")
const ErrNoKeyFound = keysError("keys: could not find any key from ssh-agent")
const ErrWrongKeyType = keysError("keys: wrong key type, only type supported is ed25519")
const ErrSshKeyRequiresPassphrase = keysError("keys: SSH Key requires a passphrase")

type keysError string

func (e keysError) Error() string {
	return string(e)
}
