Provides a resource to create a CBS snapshot.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot" "example" {
  snapshot_name = "tf-example"
  storage_id    = "disk-alc1r5sw"
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