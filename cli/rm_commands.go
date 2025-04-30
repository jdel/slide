package cli

import (
	"context"
	"errors"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *kvCommand) delete(_ *fisk.ParseContext) error {
	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()

	if c.force {
		// Actually purge the underlying message from the underlying stream
		js, err := jsClient(c.nk)
		if err != nil {
			return err
		}
		stream, err := js.Stream(ctx, "KV_"+c.bucket)
		if err != nil {
			return err
		}
		err = stream.Purge(ctx, jetstream.WithPurgeSubject("$KV."+c.bucket+"."+c.key))
		if err != nil {
			return err
		}
		return nil
	}

	kv, err := c.getBucket()
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			return nil
		}
		return err
	}

	// Otherwise soft delete
	err = kv.Delete(ctx, c.key, nil)
	if err != nil {
		return err
	}

	return nil
}

func configureDeleteCommand(app *fisk.Application) {
	c := &kvCommand{}

	rm := app.Command("rm", "Marks value as deleted").Alias("delete").Action(c.delete)
	rm.CheatFile(cheatsFs, "rm", "cheats/rm.md")
	rm.PreAction(c.parseAtBucket)
	rm.PreAction(c.getKeyPair)
	rm.Arg("key", "Name of the key[@bucket] to mark as deleted").Required().PlaceHolder("key[@bucket]").StringVar(&c.key)
	rm.Flag("purge", "Permanently delete key[@bucket] from storage").UnNegatableBoolVar(&c.force)
}
