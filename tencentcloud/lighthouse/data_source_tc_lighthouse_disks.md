Use this data source to query detailed information of lighthouse disk

Example Usage

```hcl
data "tencentcloud_lighthouse_disks" "disks" {
  disk_ids = ["lhdisk-xxxxxx"]
}
```