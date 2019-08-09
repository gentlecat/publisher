workflow "Build workflow" {
  on = "push"
  resolves = [
    "Publish the image",
    "Docker Registry",
    "Go test",
  ]
}

action "Docker Registry" {
  uses = "actions/docker/login@0c53e4a"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Build the image" {
  uses = "actions/docker/cli@0c53e4a"
  needs = ["Docker Registry"]
  args = "build -t rtsk/publisher ."
}

action "Publish the image" {
  uses = "actions/docker/cli@0c53e4a"
  needs = ["Build the image"]
  args = "push rtsk/publisher"
}

action "Go test" {
  uses = "actions/setup-go@419ae75c254126fa6ae3e3ef573ce224a919b8fe"
  runs = "go test ./... -bench ."
}
