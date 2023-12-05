Use this data source to query the list of SQL Server readonly groups.

Example Usage

```hcl
data "tencentcloud_sqlserver_dbs" "example" {
  instance_id = "mssql-ds1xhnt9"
}
```