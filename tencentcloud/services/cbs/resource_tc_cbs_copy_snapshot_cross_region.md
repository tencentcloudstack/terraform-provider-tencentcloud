Provides a resource to create a CBS copy snapshot cross region resource.

Example Usage

```hcl
resource "tencentcloud_cbs_copy_snapshot_cross_region" "example" {
  snapshot_id        = "snap-xxxxxxxx"
  destination_regions = ["ap-shanghai", "ap-chengdu"]

  snapshot_name = "my-copied-snapshot"
}
```

Import

CBS copy snapshot cross region can be imported using the composite ID, e.g.

```
$ terraform import tencentcloud_cbs_copy_snapshot_cross_region.example snap-xxxxxxxx#snap-yyyyyyyy
```
