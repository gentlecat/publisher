FROM docker.io/golang:alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c AS builder
WORKDIR /usr/src/publisher
COPY . .
RUN go build -v -o bin/publisher go.roman.zone/publisher/cmd/publisher

FROM docker.io/alpine:latest@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a AS runtime
#RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/src/publisher/bin/publisher /usr/local/bin/publisher
WORKDIR /content
CMD ["publisher", "-content", "/content", "-out", "/output"]
