Provides a resource to create a cbs snapshot_share_permission

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_share_permission" "snapshot_share_permission" {
  account_ids = ["1xxxxxx", "2xxxxxx"]
  snapshot_id = "snap-xxxxxx"
}
```

Import

cbs snapshot_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission snap-xxxxxx
```