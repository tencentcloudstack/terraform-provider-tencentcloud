Provides a resource to create a CCN instance.

Example Usage

Create a prepaid CCN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
```

Create a post-paid regional export speed limit type CCN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "OUTER_REGION_LIMIT"
}
```

Create a post-paid inter-regional rate limit type CNN

```hcl
resource "tencentcloud_ccn" "main" {
  name                 = "ci-temp-test-ccn"
  description          = "ci-temp-test-ccn-des"
  qos                  = "AG"
  charge_type          = "POSTPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
}
```

Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn.test ccn-id
```