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

func (c *kvCommand) listKeys(_ *fisk.ParseContext) error {
	kv, err := c.getBucket()
	if err != nil {
		if errors.Is(err, jetstream.ErrBucketNotFound) {
			return nil
		}
		return err
	}

	var keys []string
	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	lister, err := kv.ListKeys(ctx, nil)
	if err != nil {
		return err
	}
	for k := range lister.Keys() {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)

	fmt.Println(strings.Join(keys, "\n"))
	return nil
}

func configureListCommand(app *fisk.Application) {
	c := &kvCommand{}

	ls := app.Command("ls", "Lists keys from KV").Alias("list").Action(c.listKeys)
	ls.CheatFile(cheatsFs, "ls", "cheats/ls.md")
	ls.PreAction(c.parseAtBucket)
	ls.PreAction(c.getKeyPair)
	ls.Arg("bucket", "Optional @bucket").Default("@").PlaceHolder("@bucket").StringVar(&c.key)
}
