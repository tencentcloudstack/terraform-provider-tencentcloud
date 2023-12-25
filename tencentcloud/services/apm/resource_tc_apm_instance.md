Provides a resource to create a apm instance

~> **NOTE:** To use the field `pay_mode`, you need to contact official customer service to join the whitelist.

Example Usage

```hcl
resource "tencentcloud_apm_instance" "instance" {
  name = "terraform-test"
  description = "for terraform test"
  trace_duration = 15
  span_daily_counters = 20
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

apm instance can be imported using the id, e.g.

```
terraform import tencentcloud_apm_instance.instance instance_id
```