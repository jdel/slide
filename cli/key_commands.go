package cli

import (
	"crypto/ed25519"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/keys"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nkeys"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

type infoCommand struct {
	insecure bool
	json     bool
	seed     string
	force    bool
	nk       nkeys.KeyPair
}

func (c *infoCommand) show(_ *fisk.ParseContext) error {
	opts := options.Options

	userNkey, err := c.nk.PublicKey()
	if err != nil {
		return err
	}

	userSeed := make([]byte, 0)
	if c.insecure {
		userSeed, err = c.nk.Seed()
		if err != nil {
			return err
		}
	}

	info := struct {
		UserNkey string        `json:"userNkey,omitempty" yaml:"userNkey,omitempty"`
		UserSeed string        `json:"userSeed,omitempty" yaml:"userSeed,omitempty"`
		SshAgent string        `json:"sshAgent,omitempty" yaml:"sshAgent,omitempty"`
		SshKey   string        `json:"sshKey,omitempty" yaml:"sshKey,omitempty"`
		Servers  string        `json:"servers,omitempty" yaml:"servers,omitempty"`
		Timeout  time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	}{
		UserNkey: userNkey,
		UserSeed: string(userSeed),
		SshAgent: opts.SshAgent,
		SshKey:   opts.SshKey,
		Servers:  opts.Servers,
		Timeout:  opts.Timeout,
	}

	if c.json {
		strOpts, err := json.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Println(string(strOpts))
		return nil
	}

	strOpts, err := yaml.Marshal(info)
	if err != nil {
		return err
	}
	fmt.Println(string(strOpts))
	return nil
}

func (c *infoCommand) toSsh(_ *fisk.ParseContext) error {
	_, rawSeed, err := nkeys.DecodeSeed([]byte(c.seed))
	if err != nil {
		return err
	}

	sshKey := ed25519.NewKeyFromSeed(rawSeed)
	openSshPem, err := ssh.MarshalPrivateKey(sshKey, "")
	if err != nil {
		return err
	}

	if c.insecure {
		fmt.Println(string(pem.EncodeToMemory(openSshPem)))
		return nil
	}

	public := sshKey.Public()
	ed, ok := public.(ed25519.PublicKey)
	if !ok {
		return keys.ErrWrongKeyType
	}

	pubKey, err := ssh.NewPublicKey(ed)
	if err != nil {
		return err
	}

	fmt.Println(string(ssh.MarshalAuthorizedKey(pubKey)))
	return nil
}

func (c *infoCommand) toNkey(_ *fisk.ParseContext) error {
	userNkey, err := c.nk.PublicKey()
	if err != nil {
		return err
	}

	userSeed, err := c.nk.Seed()
	if err != nil {
		return err
	}

	_, err = os.Stat(userNkey + ".nk")
	if (err != nil && os.IsNotExist(err)) || c.force {
		os.WriteFile(userNkey+".nk", userSeed, 0600)
		return nil
	}
	return err
}

func ConfigureKeyCommands(app *fisk.Application) {
	c := &infoCommand{}
	key := app.Command("key", "Manage keys").Alias("k")

	info := key.Command("info", "Show information").Alias("i").Alias("show").Action(c.show)
	info.CheatFile(cheatsFs, "key info", "cheats/key/info.md")
	info.PreAction(c.getKeyPair)
	info.Flag("show-sensitive", "Show sensitive data").Short('S').UnNegatableBoolVar(&c.insecure)
	info.Flag("json", "Format to JSON instead of YAML").UnNegatableBoolVar(&c.json)

	toSsh := key.Command("to-ssh", "[EXPERIMENTAL] Convert nkey Seed to OpenSSH private key").Alias("c").Action(c.toSsh)
	toSsh.CheatFile(cheatsFs, "key to-ssh", "cheats/key/to-ssh.md")
	toSsh.Arg("seed", "Seed to convert").Required().StringVar(&c.seed)
	toSsh.Flag("private", "Show sensitive data").UnNegatableBoolVar(&c.insecure)

	toNkey := key.Command("to-nkey", "[EXPERIMENTAL] Convert OpenSSH private key to nkey file").Action(c.toNkey)
	toNkey.CheatFile(cheatsFs, "key to-nkey", "cheats/key/to-nkey.md")
	toNkey.PreAction(c.getKeyPair)
	toNkey.Flag("overwrite", "Overwrite existing nkey file").UnNegatableBoolVar(&c.force)
}

func (c *infoCommand) getKeyPair(_ *fisk.ParseContext) error {
	nk, err := getKeyPair()
	if err != nil {
		return err
	}
	c.nk = nk
	return nil
}
