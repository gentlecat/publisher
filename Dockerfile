FROM docker.io/golang:alpine AS builder
WORKDIR /usr/src/publisher
COPY . .
RUN go build -v -o bin/publisher go.roman.zone/publisher/cmd/publisher

FROM docker.io/alpine:latest AS runtime
#RUN apt-get update && apt-get install -y openssl ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/src/publisher/bin/publisher /usr/local/bin/publisher
WORKDIR /content
CMD ["publisher", "-content", "/content", "-out", "/output"]
