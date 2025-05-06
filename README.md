# Slide

Slide is a personal online key-value store backed by NATS.

## Foreword

In light of the [recent events and controversy](https://www.cncf.io/blog/2025/04/24/protecting-nats-and-the-integrity-of-open-source-cncfs-commitment-to-the-community/) this project is in a temporary archived state.

## Objective and motivation

Slide aims at providing a simple way to remotely store and retreive small chunks of data, using a regular ed25519 OpenSSH key for authentication.

The idea was born looking for an alternative to the now limited [Skate](https://github.com/charmbracelet/skate), operating only on a local database, and defunct [charm cloud](https://github.com/charmbracelet/charm) backend.

The main goal of slide is to provide a simple and secure way to load sensitive values in your `bashrc`, `zshrc` or any scripts without hardcoding them.

## Key takeaways

Some considerations to keep in mind:
  - Meant to store only small chunks of data, not large blobs, size limit is set by the backing NATS server.
  - Uses only a single ed25519 OpenSSH key for authentication, and salted private key encryption (when client-side encryption enabled).
  - Slide server is meant to be deployed behind a TLS termination endpoint and currently does not natively support TLS.
  - Performance is not a key factor, ease of use is at the top of the priority list.

For those reasons, Slide is a __PERSONAL__ key-value store, not meant to be used in a production environment or at large scale.

There are dedicated open source and commercial products for this kind of job.

## Usage example

Figure out the NATS Nkey associated to your SSH key:

```bash
# Using ssh-agent by default
❯ ./slide key info
userNkey: UCMXA3VWQW55NVOAV424CX77PRRNA75XT23TE3Q323GEYPBCVOMBKW55
sshAgent: /private/tmp/com.apple.launchd.4XIgOAn1G2/Listeners
servers: localhost:4222
timeout: 10s

# Using a specifig SSH key
❯ ./slide --sshkey ~/.ssh/ed_25519 key info
userNkey: UCMXA3VWQW55NVOAV424CX77PRRNA75XT23TE3Q323GEYPBCVOMBKW55
sshKey: /Users/user/.ssh/id_ed25519
servers: localhost:4222
timeout: 10s
```

Run a Slide server:

```bash
# Run the server for a single user UCMXA3VWQW55NVOAV424CX77PRRNA75XT23TE3Q323GEYPBCVOMBKW55
❯ ./slide server --user UCMXA3VWQW55NVOAV424CX77PRRNA75XT23TE3Q323GEYPBCVOMBKW55
[8654] [INF] Starting nats-server
...
[8654] [INF] Listening for client connections on 0.0.0.0:4222
[8654] [INF] Server is ready
```

Start using slide:

```bash
❯ slide set gh-token@secret GITHUB_TOKEN
❯ slide get gh-token@secret
GITHUB_TOKEN%
```

## Installation

Use a package manager:

```bash
# macOS or Linux
brew tap jdel/tap && brew install jdel/tap/slide

# Other packages options will follow
```

Or download [Binaries](releases) for Linux, macOS, FreeBSD and Windows

Or just install it with `go`:

```bash
go install github.com/jdel/slide@latest
```

## Buckets

Buckets are individual logical stores created the first time a key is inserted and are addressed using the `@bucket` syntax.

If you are mostly relying on a single bucket, you should set the `SLIDE_BUCKET` environment variable and skip the `@bucket` syntax alltogether.

```bash
❯ slide bucket ls
@secrets

❯ slide set stuff@todo "take out the trash"

❯ slide bucket ls # todo bucket has been created on demand
@secrets
@todo

❯ slide get stuff@secrets # Empty, does not exist

❯ slide get stuff@todo
take out the trash%

❯ slide ls # No default bucket set
slide: error: nats: invalid bucket name

❯ export SLIDE_BUCKET=secrets

❯ slide ls 
bank-safe-code

❯ slide get bank-safe-code
1234%
```

Each bucket is backed by a NATS KV store and an underlying stream behind the scenes.

Keep in mind that when using a SaaS provider, this amount is very likely to be limited.

## Slide server and alternatives

Slide server is nothig but a single node NATS server preconfigured for personal use with Slide.

### Self hosted NATS Server

Slide server supports taking a full NATS server configuration file if you need more elaborate configuration, but at this point, you might just run a vanilla [NATS server](https://github.com/nats-io/nats-server).

You can use Slide with your ed25519 SSH key to connect to any NATS server that supports NKeys authentication.

### SaaS prodivders

You can also use Slide with a free [Synadia Cloud](https://cloud.synadia.com) account, or other more traditional cloud providers like [Scaleway](https://www.scaleway.com).

SaaS providers will most likely run decentralized authentication instead of NKeys, which does not allow SSH keys to be used as the sole login method, a JWT will also be required.

Slide supports authentication via creds file too using `--creds` in an experimental (understand not very tested) way.

Alternatively, you could extract the ssh key from the creds file using `slide key to-ssh <SEED> --private > id_ed25519_cloud_account`, load it in your ssh-agent with `ssh-add id_ed25519_cloud_account`, then specify the JWT with `jwt` or the `SLIDE_JWT` environment variable.

## Client side encryption

If you are using a third party SaaS provider, you might want to use client side encryption with `--encrypted` or `-E` on the `set` and `get` commands.

This will perform private key encryption based on your private key and a random salt. This type of encryption is not performant on large amounts of data.

In its current form, it only works when specifying `--ssh-key` or `--creds` directly because the ssh-agent does not expose the private key, only the signature API.

## License

[MIT](https://github.com/jdel/slide/raw/main/LICENSE)
