Provides a resource to create CBS set.

Example Usage

```hcl
resource "tencentcloud_cbs_storage_set" "storage" {
        disk_count 		  = 10
        storage_name      = "mystorage"
        storage_type      = "CLOUD_SSD"
        storage_size      = 100
        availability_zone = "ap-guangzhou-3"
        project_id        = 0
        encrypt           = false
}
```