FROM docker.io/golang:alpine@sha256:d200d2dd76b03a99b52651fd1325aff7dcd35422c35c01bac56571b0507f3deb AS builder
WORKDIR /usr/src/publisher
COPY . .
RUN go build -v -o bin/publisher go.roman.zone/publisher/cmd/publisher

FROM docker.io/alpine:latest@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS runtime
#RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/src/publisher/bin/publisher /usr/local/bin/publisher
WORKDIR /content
CMD ["publisher", "-content", "/content", "-out", "/output"]
