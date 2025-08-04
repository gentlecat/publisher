# Publisher

## Usage

Check *example-content* directory for an example of what the content structure should look like.
You can build the example by running `make build-example-docker` or `make build-example`.

### Docker image build

Install Docker or compatible container runtime, then run:

```shell
docker run -v ./example-content:/content ghcr.io/gentlecat/publisher:latest
```

Built content will be in the *./example-content/out* directory.

### Build from source

First, make sure you have [Go](https://golang.org/doc/install) installed. After that, install the *publisher* itself locally:

```shell
$ go get -u go.roman.zone/publisher/cmd/publisher
```

Then run the command to generate the content:

```shell
$ publisher \
    -content "./example-content" \
    -out "./example-content/out" \
    -draft
```
