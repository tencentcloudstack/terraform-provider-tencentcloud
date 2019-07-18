resource tencentcloud_vpc vpc1 {
  name         = "ci-temp-test-vpc"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource tencentcloud_vpc vpc2 {
  name         = "ci-temp-test-vpc"
  cidr_block   = "192.168.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource tencentcloud_ccn main {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource tencentcloud_ccn_attachment attachment1 {
  ccn_id          = "${tencentcloud_ccn.main.id}"
  instance_type   = "VPC"
  instance_id     = "${tencentcloud_vpc.vpc1.id}"
  instance_region = "${var.region}"
}

resource tencentcloud_ccn_attachment attachment2 {
  ccn_id          = "${tencentcloud_ccn.main.id}"
  instance_type   = "VPC"
  instance_id     = "${tencentcloud_vpc.vpc2.id}"
  instance_region = "${var.region}"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = "${tencentcloud_ccn.main.id}"
  region          = "${var.other_region}"
  bandwidth_limit = 500
}
