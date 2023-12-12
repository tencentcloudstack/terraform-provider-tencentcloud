Use this data source to query detailed information of dbbrain mysql_process_list

Example Usage

```hcl
data "tencentcloud_dbbrain_mysql_process_list" "mysql_process_list" {
  instance_id = local.mysql_id
  product     = "mysql"
}
```