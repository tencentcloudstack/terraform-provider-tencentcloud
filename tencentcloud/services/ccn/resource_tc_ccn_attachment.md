Provides a CCN attaching resource.

Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "otheruin" {
  default = "123353"
}

variable "otherccn" {
  default = "ccn-151ssaga"
}

resource "tencentcloud_vpc" "vpc" {
  name         = "ci-temp-test-vpc"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.main.id
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
}

resource "tencentcloud_ccn_attachment" "other_account" {
  ccn_id          = var.otherccn
  instance_type   = "VPC"
  instance_id     = tencentcloud_vpc.vpc.id
  instance_region = var.region
  ccn_uin         = var.otheruin
}
```