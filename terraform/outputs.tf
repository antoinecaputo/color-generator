output "container_id" {
  description = "ID of the Docker container"
  value       = docker_container.color_generator.id
}

output "image_id" {
  description = "ID of the Docker image"
  value       = docker_image.color_generator.id
}
