run : fmt build
	go run -race main.go

fmt :
	$(info Reformatting all source files...)
	go fmt ./...

build :
	go build

serve:
	cd out && python3 -m http.server 8080
