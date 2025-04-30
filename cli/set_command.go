package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nkeys"
)

func (c *kvCommand) set(_ *fisk.ParseContext) error {
	kv, err := c.getOrCreateBucket()
	if err != nil {
		return err
	}

	// Read from stdin if no value was provided
	if c.value == "" {
		fmt.Println("Enter value. CTRL+D to end.")
		stdInBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		if stdInBytes == nil {
			return fmt.Errorf("missing value")
		}
		c.value = string(stdInBytes)
	}

	finalValue := []byte(c.value)
	if c.crypt {
		// First, use provided seed if any
		if len(c.seed) > 0 {
			seed := []byte(c.seed)
			c, err := nkeys.FromCurveSeed(seed)
			if err != nil {
				return err
			}
			recipient, err := c.PublicKey()
			if err != nil {
				return err
			}
			finalValue, err = c.Seal(finalValue, recipient)
			if err != nil {
				return err
			}
		} else {
			// recipient is overridden in Seal implementation
			finalValue, err = c.nk.Seal(finalValue, "")
			if err != nil {
				return err
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	_, err = kv.Put(ctx, c.key, finalValue)
	if err != nil {
		return err
	}
	return nil
}

func configureSetCommand(app *fisk.Application) {
	c := &kvCommand{}

	set := app.Command("set", "Stores value in bucket").Alias("put").Alias("p").Action(c.set)
	set.CheatFile(cheatsFs, "set", "cheats/set.md")
	set.PreAction(c.parseAtBucket)
	set.PreAction(c.getKeyPair)
	set.Arg("key", "Name of the key[@bucket] to write").Required().PlaceHolder("key[@bucket]").StringVar(&c.key)
	set.Arg("value", "Value of the key to write").StringVar(&c.value)
	set.Flag("encrypted", "Perform client side encryption").Short('E').UnNegatableBoolVar(&c.crypt)
	set.Flag("seed", "Curve seed for client side encryption (USE ENV VAR!)").Short('S').Envar(options.EnvPrefix + "CURVE_SEED").StringVar(&c.seed)
}
