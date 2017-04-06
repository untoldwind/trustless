# Trustless store

By default the trustless store is created at `<HOME>/.trustless_store`. With
the following layout:

| File                 | Criticality    | Description                                                       |
| -------------------- | -------------- | ----------------------------------------------------------------- |
| ring                 | VERY IMPORTANT | Contains the pgp key ring (in export format)                      |
| blocks/prefix/id     | IMPORTANT      | Contains an encrypted secret version                              |
| index/nodeId         | uncritical     | Contains an encrypted index for a node, will be recreated if lost |
| logs/nodeId          | uncritical     | Contains change log for each node, can be restored                |

Not the the ID of a block is always the hex encoded SHA-256 sum of its content. This allows a very simple integrity check:

```
for i in $(find . -type f); do sha256sum $i; done
```
(this can probably done much fancier)

## Use with gpg

Import keyring:

```
gpg --no-default-keyring --keyring trustless.sec --import ring
```
(gpg should ask for your passphrase)

```
gpg --no-default-keyring --keyring trustless.sec --decrypt <block-file> | sed -z '2q;d' | jq
```
(gpg should ask again for your passphrase)
