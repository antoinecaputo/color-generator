variable "author" {}

module "docker" {
  source = "./modules/docker"
  author = var.author
}