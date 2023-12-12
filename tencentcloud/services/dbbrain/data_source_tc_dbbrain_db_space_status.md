Use this data source to query detailed information of dbbrain db_space_status

Example Usage

```hcl
data "tencentcloud_dbbrain_db_space_status" "db_space_status" {
  instance_id = "%s"
  range_days = 7
  product = "mysql"
}
```