#example 1 (ccn)

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_dc_gateway" "ccn_main" {
  name                = "ci-cdg-ccn-test"
  network_instance_id = "${tencentcloud_ccn.main.id}"
  network_type        = "CCN"
  gateway_type        = "NORMAL"
}

resource "tencentcloud_dc_gateway_ccn_route" "route1" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "10.1.1.0/32"
}

resource "tencentcloud_dc_gateway_ccn_route" "route2" {
  dcg_id     = "${tencentcloud_dc_gateway.ccn_main.id}"
  cidr_block = "192.1.1.0/32"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_ccn_routes" "test" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_instances" "name_select" {
  name = "ci"
}

#You need to sleep for a few seconds because there is a cache on the server
data "tencentcloud_dc_gateway_instances" "id_select" {
  dcg_id = "${tencentcloud_dc_gateway.ccn_main.id}"
}

#example 2 (vpc)

resource "tencentcloud_vpc" "main" {
  name       = "ci-vpc-instance-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_dc_gateway" "vpc_main" {
  name                = "ci-cdg-vpc-test"
  network_instance_id = "${tencentcloud_vpc.main.id}"
  network_type        = "VPC"
  gateway_type        = "NAT"
}
