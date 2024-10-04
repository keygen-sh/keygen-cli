# Keygen CLI

[![CI](https://github.com/keygen-sh/keygen-cli/actions/workflows/test.yml/badge.svg)](https://github.com/keygen-sh/keygen-cli/actions)

CLI to interact with keygen.sh.

## Installation

To install the `keygen` CLI, you can run the following script. Alternatively,
you can install manually by downloading a binary and following [the install
instructions here](https://keygen.sh/docs/cli/).

```bash
curl -sSL https://get.keygen.sh/keygen/latest/install.sh | sh
```

## Commands

For all available commands and options, run `keygen --help`.

### Generate a key pair

Generate an Ed25519 public/private key pair. The private key will be used to
sign releases, and the public key will be used to verify upgrades within your
app. This key pair proves that a release was created by you, using a
cryptographic signature. Always keep your private key in a secure location.

**Never share your private key with anyone.**

```sh
keygen genkey
```

For more usage options run `keygen genkey --help`.

### Create a release

Create a new release. This command will create a new release object. The release's
`status` will be in a `DRAFT` state, unlisted until published.

```sh
keygen new \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --channel 'stable' \
  --version '1.0.0'
```

For more usage options run `keygen new --help`.

### Upload an artifact

Upload the artifact at `<path>` to a given release. A SHA-512 `checksum` will automatically
be generated for the release. In addition, When the `--signing-key` flag is provided,
the release will be signed using Ed25519ph.

```sh
keygen upload ./build/keygen_darwin_amd64 \
  --signing-key ~/.keys/keygen.key \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release '1.0.0' \
  --platform 'darwin' \
  --arch 'amd64'
```

For more usage options run `keygen upload --help`.

### Publish a release

Publish an existing release. This command will set the release's `status` to
`PUBLISHED`.

```sh
keygen publish \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release '1.0.0'
```

For more usage options run `keygen publish --help`.

### Tag a release

Tag an existing release. This command will set the release's `tag` to the
provided value. For example, tag a `latest` release.

```sh
keygen tag 'latest' \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release '1.0.0'
```

For more usage options run `keygen tag --help`.

### Untag a release

Untag an existing release. This command will set the release's `tag` to `nil`.
For example, untag a `latest` release before tagging a newer release.

```sh
keygen untag \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release 'latest'
```

For more usage options run `keygen untag --help`.

### Yank a release

Yank an existing release. Sometimes things go wrong, and this command will set
the release's `status` to `YANKED`, unlisting it.

```sh
keygen yank \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release '1.0.0'
```

For more usage options run `keygen yank --help`.

### Delete a release

Delete an existing release. Sometimes things go really wrong, and this command
will delete the release and all of its artifacts.

```sh
keygen del \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --release '1.0.0'
```

For more usage options run `keygen del --help`.

## Upgrading

To check for an upgrade to the CLI, run the following command and follow the
prompts to install. Unless `KEYGEN_NO_AUTO_UPGRADE` is set, or the `--no-auto-upgrade`
flag is passed, the CLI will automatically check for upgrades no more than
once a day.

```sh
keygen upgrade
```

For more usage options run `keygen upgrade --help`.
