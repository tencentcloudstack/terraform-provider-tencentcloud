Provides a resource to create CBS set.

-> **NOTE:** When creating encrypted disks, if `kms_key_id` is not entered, the product side will generate a key by default.

-> **NOTE:** When using CBS encrypted disk, it is necessary to add `CVM_QcsRole` role and `QcloudKMSAccessForCVMRole` strategy to the account.

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

Create 3 standard CBS storages with customize kms_key_id

```hcl
resource "tencentcloud_cbs_storage_set" "example" {
  disk_count        = 3
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  kms_key_id        = "b60b328d-7ed5-11ef-8836-5254009ad364"
  encrypt           = true
}
```

Create 3 encrypted CBS storage with default generated kms_key_id

```hcl
resource "tencentcloud_cbs_storage_set" "example" {
  disk_count        = 3
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = true
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
