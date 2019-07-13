workflow "Publish" {
  resolves = [
    "logout",
  ]
  on = "push"
}

action "login" {
  uses = "actions/docker/login@master"
  secrets = [
    "DOCKER_USERNAME",
    "DOCKER_PASSWORD",
  ]
  env = {
    DOCKER_REGISTRY_URL = "docker.pkg.github.com"
  }
}

action "publish" {
  uses = "elgohr/Publish-Docker-Github-Action@master"
  args = "docker.pkg.github.com/elgohr/concourse-sonarqube-notifier/concourse-sonarqube-notifier"
  needs = ["login"]
}

action "logout" {
  uses = "actions/docker/cli@master"
  args = "logout"
  needs = ["publish"]
}
