Use this data source to query detailed information of sqlserver query_xevent

Example Usage

```hcl
data "tencentcloud_sqlserver_query_xevent" "example" {
  instance_id = "mssql-gyg9xycl"
  event_type  = "blocked"
  start_time  = "2023-08-01 00:00:00"
  end_time    = "2023-08-10 00:00:00"
}
```