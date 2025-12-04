Provides a resource to create a APM sample config

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

resource "tencentcloud_apm_sample_config" "example" {
  instance_id    = tencentcloud_apm_instance.example.id
  sample_name    = "tf-example"
  sample_rate    = 90
  service_name   = "java-order-serive"
  operation_type = 0
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

Import

APM sample config can be imported using the instanceId#sampleName, e.g.

```
terraform import tencentcloud_apm_sample_config.example apm-jPr5iQL77#tf-example
```
