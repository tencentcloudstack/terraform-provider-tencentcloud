# Create VPC and Subnet
data "tencentcloud_vpc_instances" "example" {
  name = "guagua-ci-temp-test"
}
data "tencentcloud_vpc_subnets" "example" {
  name = "guagua-ci-temp-test"
}
#Create EIP
resource "tencentcloud_eip" "example" {
  name = "example"
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
