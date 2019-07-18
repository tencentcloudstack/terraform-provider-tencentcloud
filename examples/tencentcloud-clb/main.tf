resource "tencentcloud_security_group" "foo" {
  name = "ci-temp-test-sg"
}

resource "tencentcloud_vpc" "foo" {
  name       = "vpc-temp-test"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_clb_instance" "my_clb" {
  network_type              = "${var.network_type}"
  clb_name                  = "tf-test-clb"
  project_id                = 0
  vpc_id                    = "${tencentcloud_vpc.foo.id}"
  target_region_info_region = "ap-guangzhou"
  target_region_info_vpc_id = "${tencentcloud_vpc.foo.id}"


  security_groups = ["${tencentcloud_security_group.foo.id}"]
}