Provides a resource to enable mongodb transparent data encryption

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_transparent_data_encryption" "encryption" {
    instance_id = "cmgo-xxxxxx"
    kms_region = "ap-guangzhou"
}
```

Import

mongodb transparent data encryption can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_transparent_data_encryption.encryption ${instanceId}
```