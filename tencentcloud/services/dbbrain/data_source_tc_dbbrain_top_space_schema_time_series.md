Use this data source to query detailed information of dbbrain top_space_schema_time_series

Example Usage

```hcl
data "tencentcloud_dbbrain_top_space_schema_time_series" "top_space_schema_time_series" {
  instance_id = "%s"
  sort_by = "DataLength"
  start_date = "%s"
  end_date = "%s"
  product = "mysql"
}
```