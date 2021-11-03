# Keygen CLI

## Generate a key pair

Generate an Ed25519 a public/private key pair. The private key will be used to
sign releases, and the public key will be used to verify upgrades within your
appliaction. **Never share your private key.**

For more usage options run `keygen genkey --help`.

```
keygen genkey
```

## Publish a release

Publish a new release. This command will create a new release object, and then
upload the file at `<path>` to the release's artifact relationship.

For more usage options run `keygen releases publish --help`.

```
keygen releases publish <path>
```
