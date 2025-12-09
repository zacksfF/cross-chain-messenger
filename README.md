# cross-chain-messenger
cross-chain-messenger


The go build ./... command builds all packages but doesn't create named binaries. You need to specify the output:

```bash
go build -o relayerd ./cmd/relayerd
./relayerd
```

Or you can use go run directly:
```bash
go run cmd/relayerd/main.go
```

The same applies for the CLI:
```bash
cd ../cli
go build -o messenger-cli ./cmd/messenger-cli
./messenger-cli --help
```

