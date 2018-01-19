data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_cbs_storage" "my-storage" {
  storage_type = "cloudBasic"
  storage_size = 10
  period       = 1
  availability_zone    = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  storage_name = "${var.storage_name}"
}
