clean:
	-rm -r build
	-rm -r out

go-dependencies:
	cd main && go get -v -t ./...

go-fmt:
	$(info Reformatting all source files...)
	cd main && go fmt ./...

go-build: clean go-fmt go-dependencies
	cd main && go build -o ../build/publisher go.roman.zone/publisher/cmd/publisher

go-test: go-build
	cd main && go test ./... -bench .

build-example: go-build
	./build/publisher \
		-content "example-content" \
		-out "out"

serve:
	cd out && python3 -m http.server 8080
