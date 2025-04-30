package cli

import (
	"log"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/choria-io/fisk/units"
	"github.com/jdel/slide/logger"
	"github.com/jdel/slide/options"
	"github.com/nats-io/nats-server/v2/server"
)

var NS *server.Server

type serveCommand struct {
	host        string
	port        int
	users       []string
	multiTenant bool
	wsPort      int
	dataDir     string
	maxMem      units.Base2Bytes
	maxStorage  units.Base2Bytes
	cipher      string
	cipherKey   string
	configFile  string
}

func ConfigureServerCommand(app *fisk.Application) {
	c := &serveCommand{}
	server := app.Command("server", "Run embedded NATS Server. Global flags don't apply.").Alias("run").Action(c.run)

	server.Flag("host", "Host to bind to").Envar(options.EnvPrefix + "SERVER_HOST").Default("0.0.0.0").StringVar(&c.host)
	server.Flag("port", "Port to bind to").Envar(options.EnvPrefix + "SERVER_PORT").Default("4222").IntVar(&c.port)
	server.Flag("ws-port", "Use websocket").Envar(options.EnvPrefix + "SERVER_WS_PORT").IntVar(&c.wsPort)
	server.Flag("data-dir", "Directory for jetstream data").Envar(options.EnvPrefix + "SERVER_DATA_DIR").Default("data").StringVar(&c.dataDir)
	server.Flag("max-mem", "Maximum memory size").Envar(options.EnvPrefix + "SERVER_MAX_MEM").Default("100MB").BytesVar(&c.maxMem)
	server.Flag("max-data", "Maximum storage size").Envar(options.EnvPrefix + "SERVER_MAX_DATA").Default("100MB").BytesVar(&c.maxStorage)
	server.Flag("cipher-key", "Jetstream encryption cipher key").Envar(options.EnvPrefix + "SERVER_CIPHER_KEY").StringVar(&c.cipherKey)
	server.Flag("cipher", "Jetstream encryption mode, disabled unless a cipher-key is provided").Envar(options.EnvPrefix+"SERVER_CIPHER").Default("aes").EnumVar(&c.cipher, "aes", "chacha")
	server.Flag("multi-tenant", "Separate each user in their own account").Envar(options.EnvPrefix + "SERVER_MULTI_TENANT").UnNegatableBoolVar(&c.multiTenant)
	server.Flag("user", "User NKey").Envar(options.EnvPrefix + "SERVER_USER").StringsVar(&c.users)
	server.Flag("config", "Optional config file").Envar(options.EnvPrefix + "SERVER_CONFIG").StringVar(&c.configFile)
}

func (c *serveCommand) run(_ *fisk.ParseContext) error {
	var opts *server.Options

	if c.configFile != "" {
		var err error
		opts, err = server.ProcessConfigFile(c.configFile)
		if err != nil {
			return err
		}
	} else {
		cipher := server.NoCipher
		if c.cipherKey != "" {
			switch c.cipher {
			case "aes":
				cipher = server.AES
			case "chacha":
				cipher = server.ChaCha
			default:
				cipher = server.NoCipher
			}
		}

		wsConfig := server.WebsocketOpts{}
		if c.wsPort > 0 {
			wsConfig = server.WebsocketOpts{
				Host:  c.host,
				Port:  c.wsPort,
				NoTLS: true, // Disabled for now, use behind TLS termination endpoint
			}
		}

		users := []*server.NkeyUser{}
		accounts := []*server.Account{}
		for _, nkey := range c.users {
			logger.Log.Debug().Str("nkey", nkey).Msg("Adding user to nkey auth")
			if c.multiTenant {
				account := server.NewAccount(nkey)
				accounts = append(accounts, account)
				// Does not work here: jetstream account not registered ?
				// account.EnableJetStream(map[string]server.JetStreamAccountLimits{-1, -1, -1, -1, -1, -1, -1, false})
				users = append(users, &server.NkeyUser{
					Nkey:    nkey,
					Account: account,
				})
			} else {
				users = append(users, &server.NkeyUser{
					Nkey: nkey,
				})
			}
		}

		opts = &server.Options{
			// TLSConfig: serverTlsConfig, // Disabled for now, use WS behind TLS termination endpoint
			Port:               c.port,
			Host:               c.host,
			AllowNonTLS:        true,
			JetStream:          true,
			JetStreamCipher:    cipher,
			JetStreamKey:       c.cipherKey,
			JetStreamMaxMemory: int64(c.maxMem),
			JetStreamMaxStore:  int64(c.maxStorage),
			JetStreamDomain:    "slide",
			ServerName:         "slide",
			StoreDir:           c.dataDir,
			DontListen:         c.port == 0,
			Websocket:          wsConfig,
			Nkeys:              users,
			Accounts:           accounts,
			Debug:              strings.ToLower(options.Options.LogLevel) == "debug",
		}
	}

	ns, err := server.NewServer(opts)
	if err != nil {
		log.Fatalf("server init: %v", err)
	}

	ns.ConfigureLogger()

	go ns.Start()
	defer ns.Shutdown()

	if !ns.ReadyForConnections(options.Options.Timeout) {
		log.Fatalln("NATS server timeout")
	}

	if c.configFile == "" && c.multiTenant {
		// Not sure why but js needs to be enabled with LookupAccount
		// when server.ProcessConfigFile can just enable it
		// otherwise we get error: jetstream account not registered
		// See https://github.com/nats-io/nats-server/issues/6261
		for _, account := range opts.Accounts {
			registeredAccount, err := ns.LookupAccount(account.Name)
			if err != nil {
				return err
			}
			registeredAccount.EnableJetStream(
				map[string]server.JetStreamAccountLimits{
					"": {
						MaxMemory:            -1,
						MaxStore:             -1,
						MaxStreams:           -1,
						MaxConsumers:         -1,
						MaxAckPending:        -1,
						MemoryMaxStreamBytes: -1,
						StoreMaxStreamBytes:  -1,
						MaxBytesRequired:     false,
					},
				},
			)
		}
	}

	ns.WaitForShutdown()
	return nil
}
