gofmt :
	$(info Reformatting all source files...)
	go fmt ./...

build :
	go build

run : gofmt build
	go run -race main.go

docker : gofmt
	docker-compose up --build
