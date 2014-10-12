# go get _version_ [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rsenk330/gogetver)

Stop worrying about your Go dependencies breaking with newer versions.

_This is still very much a work in progress. Currently, the only supported VCS is git._

## Usage

Pinning dependencies for your Go applications is easy. Simply add `gogetver.com` to the beginning of your import statements:

```go
import "gogetver.com/github.com/rsenk330/gogetver"
```

If you don't specify any version information, it will continue doing the default behavior (grabs master).

To pin a version, add it to the end of the package path:

```go
import "gogetver.com/github.com/rsenk330/gogetver.v1.0"
```

### How Versions Work

You are not forced into any particular naming scheme for your versions. This package will simply look for a branch or tag name matching the version specify.

## Development

`go get` does not appear to work with ports, so to test locally, you'll need to use something like [ngrok](https://ngrok.com/):

```bash
$ ngrok 3000
```

Once `ngrok` is up and running, you'll need to start the site. [Gin](https://github.com/codegangsta/gin) is useful if you want to handle automatic reloading.

Using the hostname that was assigned to you by `ngrok`, pass start the app:

```bash
$ HOSTNAME=<HOSTNAME> gin
```

You can also run it without `gin`:

```bash
$ HOSTNAME=<HOSTNAME> go run server.go
```

## Testing

```bash
$ go test
```
