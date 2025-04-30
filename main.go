package main

import (
	"errors"
	"os"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/cli"
	"github.com/jdel/slide/logger"
	"github.com/jdel/slide/options"
)

func main() {
	cliOpts := cli.Configure()

	help := `Remote Personal Key Value CLI

	See 'slide cheat' for a quick cheatsheet of commands`

	fisk.EnableFileExpansion = false
	slide := fisk.New("slide", help).DefaultEnvars()
	slide.Author("jdel")
	slide.UsageWriter(os.Stdout)
	slide.Version(options.Version)
	slide.VersionFlag.Short('v')
	slide.HelpFlag.Short('h')
	slide.WithCheats().CheatCommand.Hidden()
	slide.PreAction(func(pc *fisk.ParseContext) error {
		logger.SetLevel(cliOpts.LogLevel)
		// Don't enforce credentials requirement on serve function
		if pc.SelectedCommand == nil || // nil if no args are passed
			(pc.SelectedCommand.Model() != nil && pc.SelectedCommand.Model().Name == "server") {
			return nil
		}

		// Enforce at least one credentials method
		if cliOpts.Creds == "" && cliOpts.SshKey == "" && cliOpts.SshAgent == "" {
			// Defaults to the default ssh-agent socket
			if agentSocket := os.Getenv("SSH_AUTH_SOCK"); agentSocket != "" {
				cliOpts.SshAgent = agentSocket
				return nil
			}
			return errors.New("must specify creds, ssh key or ssh agent socket")
		}
		return nil
	})

	slide.Flag("server", "NATS server URL").Short('s').Default("localhost:4222").StringVar(&cliOpts.Servers)
	slide.Flag("bucket", "Default bucket to use if unspecified with @bucket syntax").StringVar(&cliOpts.DefaultBucket)
	slide.Flag("stream-max-size", "Max byte size when creating new buckets (might be required by SaaS providers)").BytesVar(&cliOpts.StreamMaxBytes)
	slide.Flag("creds", "Path to NATS creds file").StringVar(&cliOpts.Creds)
	slide.Flag("ssh-key", "Path to SSH Private Key (ed25519)").StringVar(&cliOpts.SshKey)
	slide.Flag("jwt", "[EXPERIMENTAL] NATS JWT").StringVar(&cliOpts.Jwt)
	slide.Flag("timeout", "Timeout for remote operations").Default("10s").DurationVar(&cliOpts.Timeout)
	slide.Flag("ssh-key-passphrase", "SSH Key passphrase (NOT RECOMMENDED)").StringVar(&cliOpts.SshKeyPassphrase)
	slide.Flag("ssh-agent", "Use SSH Agent socket").StringVar(&cliOpts.SshAgent)
	slide.Flag("log-level", "Log level").Short('l').Default("fatal").EnumVar(&cliOpts.LogLevel, "fatal", "error", "warning", "info", "debug")

	cli.ConfigureKvCommands(slide)
	cli.ConfigureKeyCommands(slide)
	cli.ConfigureServerCommand(slide)

	slide.MustParseWithUsage(os.Args[1:])
	cli.Done()
}
