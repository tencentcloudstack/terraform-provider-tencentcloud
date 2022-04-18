resource "tencentcloud_instance" "cvm" {
  availability_zone = var.default_az
  image_id          = var.image_id
  instance_name     = var.cvm_name
}