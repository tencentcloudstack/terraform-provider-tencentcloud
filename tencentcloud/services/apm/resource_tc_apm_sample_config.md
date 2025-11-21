Provides a resource to create a APM sample config

Example Usage

```hcl
resource "tencentcloud_apm_sample_config" "example" {
  instance_id          = ""
  sample_name          = ""
  sample_rate          = 10
  service_name         = ""
  operation_name       = ""
  operation_type       = ""
  sample_config_status = 1
  tags {
    key   = "createdBy"
    value = "Terraform"
  }
}
```

Import

APM sample config can be imported using the instanceId#sampleName, e.g.

```
terraform import tencentcloud_apm_sample_config.example apm-1o8yMC47u#tf-example
```
