Provides a resource to create a CBS disk backup.

~> **NOTE:** The parameter `disk_backup_quota` in the resource `tencentcloud_cbs_storage` must be greater than 1.

Example Usage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-6"
  project_id        = 0
  encrypt           = false
  disk_backup_quota = 3

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cbs_disk_backup" "example" {
  disk_id          = tencentcloud_cbs_storage.example.id
  disk_backup_name = "tf-example"
}
```

Import

CBS disk backup can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_disk_backup.example dbp-qax6zwvr
```