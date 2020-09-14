provider "tencentcloud" {
  region = "ap-guangzhou"
}

resource "tencentcloud_image" "image_instance" {
   image_name         = var.image_imstance_name
   instance_id        = "ins-2ju245xg"
   data_disk_ids      = ["disk-gii0vtwi"]
   force_poweroff     = true
   sysprep            = false
   image_description  = "create image with instance"
}
