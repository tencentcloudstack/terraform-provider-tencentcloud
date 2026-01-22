Provides a resource to create a CBS snapshot.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot" "example" {
  storage_id    = "disk-1i9gxxi8"
  snapshot_name = "tf-example"
  disk_usage    = "DATA_DISK"
  tags = {
    createBy = "Terraform"
  }
}
```

Import

CBS snapshot can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot.example snap-3sa3f39b
```
