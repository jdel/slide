package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/nats-io/nkeys"
)

type kvCommand struct {
	bucket   string
	force    bool
	key      string
	value    string
	template string
	crypt    bool
	seed     string
	nk       nkeys.KeyPair
}

func ConfigureKvCommands(app *fisk.Application) {
	configureListCommand(app)
	configureGetCommand(app)
	configureSetCommand(app)
	configureDeleteCommand(app)
	configureBucketCommand(app)
}

// PreAction
func (c *kvCommand) getKeyPair(_ *fisk.ParseContext) error {
	nk, err := getKeyPair()
	if err != nil {
		return err
	}
	c.nk = nk
	return nil
}

// PreAction
func (c *kvCommand) parseAtBucket(pc *fisk.ParseContext) error {
	opts := options.Options
	splitKeyAtBucket := strings.Split(c.key, "@")
	if len(splitKeyAtBucket) > 1 {
		keyName := splitKeyAtBucket[0]
		bucketName := splitKeyAtBucket[1]
		if bucketName != "" {
			c.bucket = bucketName
		}
		c.key = keyName
	}
	if c.bucket == "" {
		c.bucket = opts.DefaultBucket
	}

	return nil
}

func (c *kvCommand) getBucket() (jetstream.KeyValue, error) {
	js, err := jsClient(c.nk)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	return js.KeyValue(ctx, c.bucket)
}

func (c *kvCommand) getOrCreateBucket() (jetstream.KeyValue, error) {
	js, err := jsClient(c.nk)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	kv, err := js.KeyValue(ctx, c.bucket)
	if errors.Is(err, jetstream.ErrBucketNotFound) {
		kv, err = js.CreateKeyValue(ctx, jetstream.KeyValueConfig{
			Bucket:      c.bucket,
			MaxBytes:    int64(options.Options.StreamMaxBytes),
			History:     64,
			Compression: true,
		})
		if errors.Is(err, jetstream.ErrBucketExists) {
			return nil, err
		}
	}
	return kv, err
}
