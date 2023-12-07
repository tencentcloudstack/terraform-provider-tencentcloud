Provides a resource to limit CCN bandwidth.

Example Usage

Set the upper limit of regional outbound bandwidth

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

resource "tencentcloud_ccn" "main" {
  name        = "ci-temp-test-ccn"
  description = "ci-temp-test-ccn-des"
  qos         = "AG"
}

resource "tencentcloud_ccn_bandwidth_limit" "limit1" {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  bandwidth_limit = 500
}
```

Set the upper limit between regions

```hcl
variable "other_region1" {
  default = "ap-shanghai"
}

variable "other_region2" {
  default = "ap-nanjing"
}

resource tencentcloud_ccn main {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}

resource tencentcloud_ccn_bandwidth_limit limit1 {
  ccn_id          = tencentcloud_ccn.main.id
  region          = var.other_region1
  dst_region      = var.other_region2
  bandwidth_limit = 100
}
```