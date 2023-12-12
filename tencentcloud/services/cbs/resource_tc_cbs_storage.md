Provides a resource to create a CBS.

Example Usage

```hcl
resource "tencentcloud_cbs_storage" "storage" {
  storage_name      = "mystorage"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    test = "tf"
  }
}
```

Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.storage disk-41s6jwy4
```