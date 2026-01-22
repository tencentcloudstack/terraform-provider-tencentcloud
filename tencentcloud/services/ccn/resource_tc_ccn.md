Provides a resource to create a CCN instance.

~> **NOTE:** `route_overlap_flag` currently does not support setting to `false`.

Example Usage

Create a PREPAID CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                   = "tf-example"
  description            = "desc."
  qos                    = "AG"
  charge_type            = "PREPAID"
  bandwidth_limit_type   = "INTER_REGION_LIMIT"
  instance_metering_type = "BANDWIDTH"
  route_ecmp_flag        = true
  route_overlap_flag     = true

  tags = {
    createBy = "Terraform"
  }
}
```

Create a POSTPAID regional export speed limit type CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
  route_ecmp_flag      = false
  route_overlap_flag   = true
  tags = {
    createBy = "Terraform"
  }
}
```

Create a POSTPAID inter-regional rate limit type CNN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
```

Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn.example ccn-al70jo89
```
