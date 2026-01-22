Provides a resource to create a CBS snapshot share permission

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_share_permission" "example" {
  account_ids = ["10002320****", "10002277****"]
  snapshot_id = "snap-cs5kj0j8"
}
```

Import

CBS snapshot share permission can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_share_permission.example snap-cs5kj0j8
```
