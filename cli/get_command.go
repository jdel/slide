package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
)

func (c *kvCommand) get(_ *fisk.ParseContext) error {
	kv, err := c.getBucket()
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	entry, err := kv.Get(ctx, c.key)
	if err != nil {
		if errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil
		}
		return err
	}

	finalValue := entry.Value()
	if c.crypt {
		// First, use provided seed if any
		if len(c.seed) > 0 {
			seed := []byte(c.seed)
			c, err := nkeys.FromCurveSeed(seed)
			if err != nil {
				return err
			}
			sender, err := c.PublicKey()
			if err != nil {
				return err
			}
			finalValue, err = c.Open(finalValue, sender)
			if err != nil {
				return err
			}
		} else {
			// sender is overridden in Open implementation
			finalValue, err = c.nk.Open(finalValue, "")
			if err != nil {
				return err
			}
		}
	}
	fmt.Print(string(finalValue))
	return nil
}

func configureGetCommand(app *fisk.Application) {
	c := &kvCommand{}

	get := app.Command("get", "Gets value from bucket").Alias("g").Action(c.get)
	get.CheatFile(cheatsFs, "get", "cheats/get.md")
	get.PreAction(c.parseAtBucket)
	get.PreAction(c.getKeyPair)
	get.Arg("key", "Name of the key[@bucket] to fetch").Required().PlaceHolder("key[@bucket]").StringVar(&c.key)
	get.Flag("encrypted", "Perform client side decryption").Short('E').UnNegatableBoolVar(&c.crypt)
	get.Flag("seed", "Curve seed for decryption (USE ENV VAR!)").Short('S').Envar(options.EnvPrefix + "CURVE_SEED").StringVar(&c.seed)
}
