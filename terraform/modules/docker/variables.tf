variable "author" {}

variable "container_name" {
  description = "Value of the name for the Docker container"
  type        = string
  default     = "ColorGeneratorContainer"
}

variable "volume_name" {
  description = "Value of the name for the Docker volume"
  type        = string
  default     = "ColorGeneratorVolume"
}

