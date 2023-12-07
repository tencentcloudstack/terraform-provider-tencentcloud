Use this data source to query detailed information of CBS snapshots.

Example Usage

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_id        = "snap-f3io7adt"
  result_output_file = "mytestpath"
}
```