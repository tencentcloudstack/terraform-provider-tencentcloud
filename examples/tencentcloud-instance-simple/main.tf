resource "tencentcloud_instance" "my-server" {
  image_id          = var.image_id
  availability_zone = var.availability_zone
}

output "instance_id" {
  value = tencentcloud_instance.my-server.id
}
