Provides a resource to create a Cynosdb cluster transparent encrypt

~> **NOTE:** Once activated, it cannot be deactivated.

~> **NOTE:** If you have not enabled the KMS service or authorized the KMS key before, you will need to enable the KMS service and then authorize the KMS key in order to complete the corresponding enabling or authorization operations and unlock the subsequent settings for data encryption.

Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster_transparent_encrypt" "example" {
  cluster_id                = cynosdbmysql-bu6hlulf
  key_id                    = "f063c18b-xxxx-xxxx-xxxx-525400d3a886"
  key_region                = "ap-guangzhou"
  key_type                  = "custom"
  is_open_global_encryption = false
}
```

Import

Cynosdb cluster transparent encrypt can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_cluster_transparent_encrypt.example cynosdbmysql-bu6hlulf
```