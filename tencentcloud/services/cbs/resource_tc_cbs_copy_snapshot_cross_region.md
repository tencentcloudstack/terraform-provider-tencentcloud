Provides a resource to create a CBS copy snapshot cross region resource.

Example Usage

```hcl
resource "tencentcloud_cbs_copy_snapshot_cross_region" "example" {
  snapshot_id        = "snap-07ttd84z"
  destination_region = "ap-beijing"
  snapshot_name      = "tf-example-copy-snapshot"
}
```
