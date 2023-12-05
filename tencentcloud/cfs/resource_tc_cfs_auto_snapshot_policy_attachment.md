Provides a resource to create a cfs auto_snapshot_policy_attachment

Example Usage

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy_attachment" "auto_snapshot_policy_attachment" {
  auto_snapshot_policy_id = "asp-basic"
  file_system_ids         = "cfs-4xzkct19,cfs-iobiaxtj"
}
```

Import

cfs auto_snapshot_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy_attachment.auto_snapshot_policy_attachment auto_snapshot_policy_id#file_system_ids
```