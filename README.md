# Publisher [![Travis CI](https://img.shields.io/travis/gentlecat/go-bomb.svg?style=flat-square)](https://travis-ci.org/gentlecat/publisher)

## Installation

First, make sure you have [Go](https://golang.org/doc/install) installed. After that, install *publisher* locally:

```bash
$ go get -u go.roman.zone/publisher
```

## Usage

```bash
$ publisher \
    -content "~/fancy_website/content" \
    -out "~/fancy_website/public" \
    -prod
```

Check *example-content* directory for an example of how content structure should look like.
