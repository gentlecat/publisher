workflow "Build workflow" {
  on = "push"
  resolves = [
    "Publish the image",
    "Docker Registry",
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
