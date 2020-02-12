# Removes all the build directories
clean :
	-rm -r build
	-rm -r out

test : build-example
	go test ./... -bench .

run : fmt
	go run -race main.go

fmt :
	$(info Reformatting all source files...)
	go fmt ./...

build : clean fmt
	go build -o ./build/publisher go.roman.zone/publisher/cmd/publisher

build-example : build
	./build/publisher \
		-content "example-content" \
		-out "out"

serve:
	cd out && python3 -m http.server 8080
