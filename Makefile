clean :
	-rm -r build
	-rm -r out

run : fmt
	go run -race main.go

fmt :
	$(info Reformatting all source files...)
	go fmt ./...

build : clean
	-rm -r ./build
	go build -o ./build/publisher

serve:
	cd out && python3 -m http.server 8080
