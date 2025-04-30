package cli

import (
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/jdel/slide/keys"
	"github.com/jdel/slide/logger"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
)

//go:embed cheats
var cheatsFs embed.FS

var (
	start time.Time
)

func Configure() *options.CliOptions {
	start = time.Now()
	options.Options = &options.CliOptions{}
	return options.Options
}

func Done() {
	logger.Log.Debug().Dur("duration", time.Since(start)).Msg("Complete")
}

func getKeyPair() (nkeys.KeyPair, error) {
	opts := options.Options
	if opts.Creds != "" {
		logger.Log.Debug().Str("file", opts.Creds).Msg("Getting key from NATS creds file")
		return keys.FromCreds(opts.Creds)
	} else if opts.SshKey != "" {
		logger.Log.Debug().Str("file", opts.SshKey).Msg("Getting key from ssh private key file")
		return keys.FromFile(opts.SshKey, opts.SshKeyPassphrase)
	} else if opts.SshAgent != "" {
		logger.Log.Debug().Str("agent", opts.SshAgent).Msg("Getting key from ssh-agent")
		// TODO: Maybe handle selecting which key somehow based on key comment ?
		kp, err := keys.GetKeyPairsFromAgent(opts.SshAgent)
		if err != nil {
			return nil, err
		}
		return kp[0], nil
	}
	return nil, errors.New("no key specified")
}

func natsClient(nk nkeys.KeyPair) (*nats.Conn, error) {
	opts := options.Options

	userNkey, err := nk.PublicKey()
	if err != nil {
		return nil, err
	}

	switch {
	case opts.Creds != "":
		logger.Log.Debug().
			Str("nkey", userNkey).
			Str("server", opts.Servers).
			Msg("Connecting with Creds")
		return nats.Connect(opts.Servers, nats.UserCredentials(opts.Creds))

	case opts.Jwt != "":
		logger.Log.Debug().
			Str("nkey", userNkey).
			Str("jwt", opts.Jwt).
			Str("server", opts.Servers).
			Msg("Connecting with SSH Key + JWT")
		return nats.Connect(opts.Servers, nats.UserJWT(
			func() (string, error) {
				return opts.Jwt, nil
			}, func(nonce []byte) ([]byte, error) {
				return nk.Sign(nonce)
			}))

	case opts.SshKey != "" || opts.SshAgent != "":
		logger.Log.Debug().
			Str("nkey", userNkey).
			Str("server", opts.Servers).
			Msg("Connecting with SSH Key")
		return nats.Connect(opts.Servers, nats.Nkey(
			userNkey,
			func(nonce []byte) ([]byte, error) {
				return nk.Sign(nonce)
			},
		))

	default:
		return nil, fmt.Errorf("no suitable authentication method found")
	}
}

func jsClient(nk nkeys.KeyPair) (jetstream.JetStream, error) {
	nc, err := natsClient(nk)
	if err != nil {
		return nil, err
	}

	return jetstream.New(nc)
}
