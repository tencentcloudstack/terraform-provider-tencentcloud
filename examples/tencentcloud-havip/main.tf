# Create VPC and Subnet
resource "tencentcloud_vpc" "example" {
  name       = "example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "example" {
  name              = "example"
  availability_zone = "${var.availability_zone}"
  vpc_id            = "${tencentcloud_vpc.example.id}"
  cidr_block        = "10.0.0.0/24"
  is_multicast      = false
}

resource "tencentcloud_ha_vip" "example" {
  name      = "example"
  vpc_id    = "${data.tencentcloud_vpc_instances.example.instance_list.0.vpc_id}"
  subnet_id = "${data.tencentcloud_vpc_subnets.example.instance_list.0.subnet_id}"
  vip       = "10.0.20.5"

}
resource "tencentcloud_ha_vip_eip_attachment" "example" {
  havip_id   = "${tencentcloud_ha_vip.example.id}"
  address_ip = "${tencentcloud_eip.example.public_ip}"
}

data "tencentcloud_ha_vips" "example" {
  id = "${tencentcloud_ha_vip.example.id}"
}
data "tencentcloud_ha_vip_eip_attachments" "example" {
  havip_id = "${tencentcloud_ha_vip_eip_attachment.example.havip_id}"
}
