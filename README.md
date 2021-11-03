# Keygen CLI

## Generate a key pair

Generate an Ed25519 public/private key pair. The private key will be used to
sign releases, and the public key will be used to verify upgrades within your
app. **Never share your private key.**

For more usage options run `keygen genkey --help`.

```sh
keygen genkey
```

## Publish a release

Publish a new release. This command will create a new release object, and then
upload the file at `<path>` to the release's artifact relationship. When the
`--signing-key` flag is provided, the release will be signed using Ed25519ph.
In addition, a SHA-512 checksum will be generated for the release.

For more usage options run `keygen releases publish --help`.

```sh
keygen releases publish dist/App-1-0-0.zip \
  --signing-key ./keygen.key \
  --account '1fddcec8-8dd3-4d8d-9b16-215cac0f9b52' \
  --product '2313b7e7-1ea6-4a01-901e-2931de6bb1e2' \
  --token 'prod-xxx' \
  --platform 'darwin' \
  --version '1.0.0'
```
