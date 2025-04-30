package cli

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats.go/jetstream"
)

func (c *kvCommand) listBuckets(_ *fisk.ParseContext) error {
	js, err := jsClient(c.nk)
	if err != nil {
		return err
	}

	var buckets []string
	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	lister := js.KeyValueStoreNames(ctx)
	for n := range lister.Name() {
		buckets = append(buckets, "@"+n)
	}
	if len(buckets) == 0 {
		return nil
	}
	slices.Sort(buckets)

	fmt.Println(strings.Join(buckets, "\n"))
	return nil
}

func (c *kvCommand) deleteBucket(_ *fisk.ParseContext) error {
	if c.bucket[0:1] == "@" {
		c.bucket = c.bucket[1:]
	}

	if !c.force {
		fmt.Printf("Bucket `%s` and all associated data will be deleted, rerun with --yes to confirm deletion\n", c.bucket)
		return nil
	}

	js, err := jsClient(c.nk)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	err = js.DeleteKeyValue(ctx, c.bucket)
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func configureBucketCommand(app *fisk.Application) {
	c := &kvCommand{}

	bucket := app.Command("bucket", "Make operations on buckets").Alias("db")
	bucket.CheatFile(cheatsFs, "bucket", "cheats/bucket.md")
	bucket.PreAction(c.getKeyPair)

	_ = bucket.Command("ls", "Lists available buckets").Alias("list").Action(c.listBuckets)

	delete := bucket.Command("delete", "Delete bucket").Alias("rm").Action(c.deleteBucket)
	delete.Arg("bucket", "Bucket name to delete").Required().StringVar(&c.bucket)
	delete.Flag("yes", "Confirm bucket deletion").Short('y').UnNegatableBoolVar(&c.force)
}
