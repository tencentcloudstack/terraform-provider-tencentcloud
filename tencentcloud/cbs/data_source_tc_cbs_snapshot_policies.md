Use this data source to query detailed information of CBS snapshot policies.

Example Usage

```hcl
data "tencentcloud_cbs_snapshot_policies" "policies" {
  snapshot_policy_id   = "snap-f3io7adt"
  snapshot_policy_name = "test"
}
```