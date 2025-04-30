package cli

import (
	"cmp"
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"text/template"

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

	type templatedKv struct {
		Key    string
		Value  string
		Bucket string
	}

	var keys []templatedKv

	ctx, cancel := context.WithTimeout(context.Background(), options.Options.Timeout)
	defer cancel()
	lister, err := kv.ListKeys(ctx, nil)
	if err != nil {
		return err
	}

	for k := range lister.Keys() {
		if strings.Contains(c.template, ".Value") {
			e, err := kv.Get(ctx, k)
			if err != nil {
				return err
			}
			keys = append(keys, templatedKv{
				Key:    k,
				Value:  string(e.Value()),
				Bucket: c.bucket,
			})
		} else {
			keys = append(keys, templatedKv{
				Key:    k,
				Bucket: c.bucket,
			})
		}
	}
	if len(keys) == 0 {
		return nil
	}

	slices.SortFunc(keys, func(a, b templatedKv) int {
		return cmp.Compare(a.Key, b.Key)
	})

	tmpl, err := template.New("ls").Parse(`{{- range . }}` + c.template + "\n" + `{{ end -}}`)
	if err != nil {
		return err
	}

	tmpl.Execute(os.Stdout, keys)
	return nil
}

func configureListCommand(app *fisk.Application) {
	c := &kvCommand{}

	ls := app.Command("ls", "Lists keys from KV").Alias("list").Action(c.listKeys)
	ls.CheatFile(cheatsFs, "ls", "cheats/ls.md")
	ls.PreAction(c.parseAtBucket)
	ls.PreAction(c.getKeyPair)
	ls.Arg("bucket", "Optional @bucket").Default("@").PlaceHolder("@bucket").StringVar(&c.key)
	ls.Flag("template", "Templated output").Default("{{ .Key }}").StringVar(&c.template)
}
