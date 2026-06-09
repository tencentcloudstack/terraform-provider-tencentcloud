Provides a resource to manage CDH instance.

~> **NOTE:** CHD instance not supported delete, please contact the work order for processing

Example Usage

```hcl
resource "tencentcloud_cdh_instance" "example" {
  availability_zone  = "ap-guangzhou-6"
  host_type          = "HC20"
  charge_type        = "PREPAID"
  prepaid_period     = 1
  host_name          = "tf-example"
  prepaid_renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}
```

Import

CDH instance can be imported using the id, e.g.

```
terraform import tencentcloud_cdh_instance.example host-d6s7i5q4
```