resource "tencentcloud_vpc" "tfer--ci-temp-test-updated_vpc-8beinccp" {
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = "false"
  name         = "ci-temp-test-updated"

  tags = {
    test = "test"
  }
}

resource "tencentcloud_vpc" "tfer--test-import-1_vpc-m3k357gn" {
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["183.60.82.98", "183.60.83.19"]
  is_multicast = "false"
  name         = "test-import-1"
}
