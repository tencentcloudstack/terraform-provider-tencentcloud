Use this data source to query detailed information of dbbrain sqlFilters

Example Usage

```hcl
resource "tencentcloud_dbbrain_sql_filter" "sql_filter" {
  instance_id = "mysql_ins_id"
  session_token {
    user = "user"
	password = "password"
  }
  sql_type = "SELECT"
  filter_key = "test"
  max_concurrency = 10
  duration = 3600
}

data "tencentcloud_dbbrain_sql_filters" "sql_filters" {
  instance_id = "mysql_ins_id"
  filter_ids = [tencentcloud_dbbrain_sql_filter.sql_filter.filter_id]
  }
```