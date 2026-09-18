Provides a resource to manage MongoDB instance CPU elastic scaling configuration.

~> **NOTE:** This resource is used to enable or disable CPU elastic scaling for MongoDB instances. When the resource is destroyed, CPU elastic scaling will be disabled automatically.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_cpu_scale" "example" {
  instance_id = "cmgo-xxxxxxxx"
  extra_cpu   = 2
}
```

Import

MongoDB instance CPU elastic scaling configuration can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_cpu_scale.example cmgo-xxxxxxxx
```