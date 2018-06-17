[![Build Status](https://travis-ci.org/cotox/go-toxcore-c.svg?branch=master)](https://travis-ci.org/cotox/go-toxcore-c)
[![GoDoc](https://godoc.org/github.com/cotox/go-toxcore-c?status.svg)](https://godoc.org/github.com/cotox/go-toxcore-c)

# go-toxcore-c

The golang bindings for libtoxcore

## Installation

```bash
go get -v gopkg.in/cotox/go-toxcore-c.v2
```

## Examples

```golang
import "gopkg.in/cotox/go-toxcore-c.v2"

// use custom options
opt := tox.NewToxOptions()
t := tox.NewTox(opt)
av := tox.NewToxAv(t)

// use default options
t := tox.NewTox(nil)
av := tox.NewToxAv(t)
```

## Tests

```bash
go test -v -covermode count
```

## Contributing

1. Fork it
2. Create your feature branch (``git checkout -b my-new-feature``)
3. Commit your changes (``git commit -am 'Add some feature'``)
4. Push to the branch (``git push origin my-new-feature``)
5. Create new Pull Request
