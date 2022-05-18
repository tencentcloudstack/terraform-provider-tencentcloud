resource "tencentcloud_instance" "cvm" {
  availability_zone = var.default_az
  image_id          = var.image_id
  instance_name     = var.cvm_name
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}

data "tencentcloud_instance_types" "default" {
  filter {
    name   = "zone"
    values = ["ap-guangzhou-3"]
  }
  cpu_core_count = 2
  exclude_sold_out = true
}

resource "tencentcloud_instance" "foo" {
  instance_name     = "keep-test-image-cvm"
  availability_zone = "ap-guangzhou-3"
  image_id          = data.tencentcloud_images.default.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  vpc_id            = "vpc-68vi2d3h"
  subnet_id         = "subnet-ob6clqwk"
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt = false
  }
}