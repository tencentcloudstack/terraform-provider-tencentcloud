resource "tencentcloud_dcx"  "bgp_main"
 {
    bandwidth = 900
    dc_id = "${var.dc_id}"
    dcg_id = "${var.dcg_id}"
    name = "bgp_main"
    network_type = "VPC"
    route_type = "BGP"
    vlan = 100
    vpc_id = "${var.vpc_id}"
}

resource "tencentcloud_dcx"  "static_main"
 {
    bandwidth = 900
    dc_id = "${var.dc_id}"
    dcg_id = "${var.dcg_id}"
    name = "static_main"
    network_type = "VPC"
    route_type = "STATIC"
    vlan = 301
    vpc_id = "${var.vpc_id}"
	route_filter_prefixes =["10.10.10.101/32"]
	tencent_address = "100.93.46.1/30"
	customer_address = "100.93.46.2/30"
}

data "tencentcloud_dcx_instances" "name_select"{
    name = "main"
}

data "tencentcloud_dcx_instances"  "id" {
    dcx_id = "${tencentcloud_dcx.static_main.id}"
}

data "tencentcloud_dc_instances"  "name" {
    name ="x"
}
data "tencentcloud_dc_instances"  "id" {
    dc_id="${var.dc_id}"
}



