Use this data source to query detailed information of Wedata list lineage

Example Usage

```hcl
data "tencentcloud_wedata_list_lineage" "example" {
  resource_unique_id = "fM8OgzE-AM2h4aaJmdXoPg"
  resource_type      = "TABLE"
  direction          = "INPUT"
  platform           = "WEDATA"
}
```
