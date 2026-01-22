Provides a resource to manage MongoDB instance SSL configuration.

~> **NOTE:** This resource is used to enable or disable SSL for MongoDB instances. When the resource is destroyed, SSL will be disabled automatically.

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_ssl" "example" {
  instance_id = "cmgo-xxxxxxxx"
  enable      = true
}
```

Import

MongoDB instance SSL configuration can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_instance_ssl.example cmgo-xxxxxxxx
```
