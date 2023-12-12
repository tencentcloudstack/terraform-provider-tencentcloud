Provides a resource to create a CBS snapshot.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot" "snapshot" {
  snapshot_name = "unnamed"
  storage_id    = "disk-kdt0sq6m"
}
```

Import

CBS snapshot can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot.snapshot snap-3sa3f39b
```