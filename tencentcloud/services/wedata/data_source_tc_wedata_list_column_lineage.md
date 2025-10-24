Use this data source to query detailed information of WeData list column lineage

Example Usage

```hcl
data "tencentcloud_wedata_list_column_lineage" "example" {
  table_unique_id = "B_CRyO4-3rMvNFPH_7aTaw"
  direction       = "INPUT"
  column_name     = "example_column"
  platform        = "WEDATA"
}
```
