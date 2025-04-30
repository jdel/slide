package options

import (
	"time"

	"github.com/choria-io/fisk/units"
)

const EnvPrefix = "SLIDE_"

// Version is the current application version.
// This variable is populated when building the binary with:
// -ldflags "-X github.com/jdel/slide/options.Version=${VERSION}"
var Version string

type CliOptions struct {
	Servers          string           // `json:"servers" yaml:"servers"`
	SshKeyPassphrase string           // `json:"-" yaml:"-"`
	Creds            string           // `json:"creds" yaml:"creds"`
	SshKey           string           // `json:"sshKey" yaml:"sshKey"`
	SshAgent         string           // `json:"sshAgent" yaml:"sshAgent"`
	LogLevel         string           // `json:"logLevel" yaml:"logLevel"`
	DefaultBucket    string           // `json:"defaultBucket" yaml:"defaultBucket"`
	StreamMaxBytes   units.Base2Bytes // `json:"defaultBucket" yaml:"defaultBucket"`
	Jwt              string           // `json:"jwt" yaml:"jwt"`
	Timeout          time.Duration    // `json:"timeout" yaml:"timeout"`
}

var Options *CliOptions
