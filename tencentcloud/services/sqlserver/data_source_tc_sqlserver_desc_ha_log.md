Use this data source to query detailed information of sqlserver desc_ha_log

Example Usage

```hcl
data "tencentcloud_sqlserver_desc_ha_log" "desc_ha_log" {
  instance_id = "mssql-jdk2pwld"
  start_time  = "2023-12-01 00:00:00"
  end_time    = "2023-12-15 00:00:00"
  switch_type = 1
}
```
