#!/usr/bin/env bash

gofmt -w ../.
go run -race main.go
