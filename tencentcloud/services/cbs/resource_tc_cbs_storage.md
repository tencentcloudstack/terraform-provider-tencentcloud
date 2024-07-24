Provides a resource to create a CBS storage.

Example Usage

Create a standard CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    createBy = "terraform"
  }
}
```

Create a dedicated cluster CBS storage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name         = "tf-example"
  storage_type         = "CLOUD_SSD"
  storage_size         = 100
  availability_zone    = "ap-guangzhou-4"
  dedicated_cluster_id = "cluster-262n63e8"
  charge_type          = "DEDICATED_CLUSTER_PAID"
  project_id           = 0
  encrypt              = false

  tags = {
    createBy = "terraform"
  }
}
```

Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.example disk-41s6jwy4
```
