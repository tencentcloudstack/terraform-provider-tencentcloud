Provides a resource to manage CDH instance.

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HC20"
  charge_type = "PREPAID"
  prepaid_period = 1
  host_name = "test"
  prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
```

Import

CDH instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cdh_instance.foo host-d6s7i5q4
```