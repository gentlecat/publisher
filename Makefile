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

build-example-docker-local: clean build-container
	docker run \
		-v ./example-content/:/content \
		-v ./out/:/output \
		publisher \
		publisher -content /content -out /output

build-example-docker-local-draft: clean build-container
	docker run \
		-v ./example-content/:/content \
		-v ./out/:/output \
		publisher \
		publisher -content /content -out /output -draft

build-example-docker: clean
	docker run \
		-v ./example-content/:/content \
		-v ./out/:/output \
		ghcr.io/gentlecat/publisher:latest \
		publisher -content /content -out /output

build-example-docker-draft: clean
	docker run \
		-v ./example-content/:/content \
		-v ./out/:/output \
		ghcr.io/gentlecat/publisher:latest \
		publisher -content /content -out /output -draft

serve:
	cd ./out && python3 -m http.server 8080
