Provides a resource to create a lighthouse snapshot

Example Usage

```hcl
resource "tencentcloud_lighthouse_snapshot" "snapshot" {
  instance_id = "lhins-acd1234"
  snapshot_name = "snap_20200903"
}
```