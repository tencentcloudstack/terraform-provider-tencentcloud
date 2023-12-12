Provides a resource to create a kms cloud_resource_attachment

Example Usage

```hcl
resource "tencentcloud_kms_cloud_resource_attachment" "example" {
  key_id      = "72688f39-1fe8-11ee-9f1a-525400cf25a4"
  product_id  = "mysql"
  resource_id = "cdb-fitq5t9h"
}
```

Import

kms cloud_resource_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_kms_cloud_resource_attachment.example 72688f39-1fe8-11ee-9f1a-525400cf25a4#mysql#cdb-fitq5t9h
```