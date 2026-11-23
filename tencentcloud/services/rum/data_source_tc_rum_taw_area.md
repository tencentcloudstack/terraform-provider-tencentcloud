Use this data source to query detailed information of rum taw area

Example Usage

```hcl
data "tencentcloud_rum_taw_area" "example" {
  area_ids      = [1]
  area_keys     = ["ap-guangzhou"]
  area_statuses = [1]
}
```
