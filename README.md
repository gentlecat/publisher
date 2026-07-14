# Publisher

## Usage

Check *example-content* directory for an example of what the content structure should look like.
You can build the example by running `make build-example-docker` or `make build-example`.

### Content structure

Each story is a *bundle*: a subdirectory of `stories/` holding the story's
content, metadata and static assets (such as images).

```
stories/
  example-story/
    index.md        # the story content, in Markdown
    metadata.json   # the story metadata (see below)
    image.jpg       # an asset, referenced from index.md as /example-story/image.jpg
```

The story `example-story` is built as `example-story.html` (served at
`/example-story`), with its assets copied into a sibling `/example-story/`
directory. Reference assets from the content with absolute paths that include
the story name, for example `/example-story/image.jpg`. Directories without a
`metadata.json` (for example, one used to group archived stories) are ignored.

### Story status

Each story can declare a `status` in its metadata:

- `published` (the default): appears everywhere, including the index, RSS feed and sitemap.
- `unlisted`: gets its own page but stays out of the index, RSS feed and sitemap. Reachable only by a direct link.
- `draft`: built only when the generator runs with `-draft` and skipped in production.

```json
{
  "title": "Work in progress",
  "date": "2000-Jan-12",
  "status": "draft"
}
```

The legacy `"draft": true` flag is still honoured (equivalent to `"status": "draft"`), but new content should use `status`.

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
