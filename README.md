## go-vet

go-vet is a set of static analyzers to use with "go vet". These are used on the golang edge-cloud repos to enforce certain coding standards and avoid potentional hard to find bugs. This tool does not replace or overlap with other common static analyzers like [staticheck](https://github.com/dominikh/go-tools) or [gosec](https://github.com/securego/gosec).

Currently, this checks for:
- shadowed variables
- "bad funcs" which can cause problems and should be avoided

## Install

Go 1.16+:
```
go install https://github.com/edgexr/go-vet@latest
```
Older go (run outside of your project directories):
```
go get -u https://github.com/edgexr/go-vet
```

## Usage:

The binary is a plugin for the standard go tool. From your project root:
```
go vet -vettool=$(which go-vet) ./...
```
