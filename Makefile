clean :
	-rm -r build
	-rm -r out

fmt :
	$(info Reformatting all source files...)
	go fmt ./...

build : clean fmt
	go build -o ./build/publisher go.roman.zone/publisher/cmd/publisher

build-example : build
	./build/publisher \
		-content "example-content" \
		-out "out"

test : build-example
	go test ./... -bench .

serve:
	cd out && python3 -m http.server 8080
