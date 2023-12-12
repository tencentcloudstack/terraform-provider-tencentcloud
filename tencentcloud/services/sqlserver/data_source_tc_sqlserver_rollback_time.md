Use this data source to query detailed information of sqlserver rollback_time

Example Usage

```hcl
data "tencentcloud_sqlserver_rollback_time" "example" {
  instance_id = "mssql-qelbzgwf"
  dbs         = ["keep_pubsub_db"]
}
```