Provides a resource to create a APM instance

~> **NOTE:** To use the field `pay_mode`, you need to contact official customer service to join the whitelist.

Example Usage

```hcl
resource "tencentcloud_apm_instance" "example" {
  name                = "tf-example"
  description         = "desc."
  trace_duration      = 7
  span_daily_counters = 0
  tags = {
    createdBy = "Terraform"
  }
}
```

Import

APM instance can be imported using the id, e.g.

```
terraform import tencentcloud_apm_instance.example apm-IMVrxXl1K
```
