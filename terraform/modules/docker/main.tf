terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.1"
    }
  }
}

# Providers configuration
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Infrastructure components definition
resource "docker_image" "color_generator" {
  name         = "${var.author}/color-generator:v0.0.2"
  build {
    context = "../"
    dockerfile = "../Dockerfile"
    label = {
      "author" = "Antoine"
    }
  }
}

resource "docker_container" "color_generator" {
  image = docker_image.color_generator.name
  name  = var.container_name
  ports {
    internal = 8080
    external = 8080
  }
  mounts {
    target = "/app/data"
    source = "/Users/antoine/Documents/Development/GolandProjects/color-generator/app/data"
    type   = "bind"
  }
  restart = "always"
}