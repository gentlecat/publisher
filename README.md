# Publisher

## Installation

First, make sure you have [Go](https://golang.org/doc/install) installed. After that, install the *publisher* itself locally:

```bash
$ go get -u go.roman.zone/publisher/cmd/publisher
```

## Usage

```bash
$ publisher \
    -content "~/fancy_website/content" \
    -out "~/fancy_website/public" \
    -prod
```

Check *example-content* directory for an example of what content structure should look like. You can build the example by running `make build-example`.
