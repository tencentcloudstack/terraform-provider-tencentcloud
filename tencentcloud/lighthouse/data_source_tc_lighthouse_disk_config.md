Use this data source to query detailed information of lighthouse disk_config

Example Usage

```hcl
data "tencentcloud_lighthouse_disk_config" "disk_config" {
  filters {
	name = "zone"
	values = ["ap-guangzhou-3"]
  }
}
```