Provides a resource to create a CCN instance.

Example Usage

Create a prepaid CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  route_ecmp_flag      = true
  route_overlap_flag   = true
  tags = {
    createBy = "terraform"
  }
}
```

Create a post-paid regional export speed limit type CCN

```hcl
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
  route_ecmp_flag      = false
  route_overlap_flag   = false
  tags = {
    createBy = "terraform"
  }
}
```

Create a post-paid inter-regional rate limit type CNN

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
