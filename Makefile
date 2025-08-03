clean:
	-rm -r build
	-rm -r out

go-dependencies:
	go get -v -t ./...

go-fmt:
	$(info Reformatting all source files...)
	go fmt ./...

go-build: clean go-fmt go-dependencies
	go build -o build/publisher go.roman.zone/publisher/cmd/publisher

go-test: go-build
	go test ./... -bench .

build-example: go-build
	./build/publisher \
		-content "example-content" \
		-out "out"

build-container:
	docker build -t publisher .

serve:
	cd out && python3 -m http.server 8080
