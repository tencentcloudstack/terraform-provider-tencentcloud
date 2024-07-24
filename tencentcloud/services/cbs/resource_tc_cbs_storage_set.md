Provides a resource to create CBS set.

Example Usage

Create 3 standard CBS storages

```hcl
resource "tencentcloud_cbs_storage_set" "example" {
  disk_count        = 3
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false
}
```

Create 3 dedicated cluster CBS storages

```hcl
resource "tencentcloud_cbs_storage_set" "example" {
  disk_count           = 3
  storage_name         = "tf-example"
  storage_type         = "CLOUD_SSD"
  storage_size         = 100
  availability_zone    = "ap-guangzhou-4"
  dedicated_cluster_id = "cluster-262n63e8"
  charge_type          = "DEDICATED_CLUSTER_PAID"
  project_id           = 0
  encrypt              = false
}
```
