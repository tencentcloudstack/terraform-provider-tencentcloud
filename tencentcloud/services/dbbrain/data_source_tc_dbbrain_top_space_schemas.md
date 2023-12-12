Use this data source to query detailed information of dbbrain top_space_schemas

Example Usage

```hcl
data "tencentcloud_dbbrain_top_space_schemas" "top_space_schemas" {
  instance_id = "%s"
  sort_by = "DataLength"
  product = "mysql"
}
```