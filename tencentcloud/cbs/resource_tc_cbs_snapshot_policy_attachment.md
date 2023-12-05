Provides a CBS snapshot policy attachment resource.

Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_policy_attachment" "foo" {
  storage_id         = tencentcloud_cbs_storage.foo.id
  snapshot_policy_id = tencentcloud_cbs_snapshot_policy.policy.id
}
```