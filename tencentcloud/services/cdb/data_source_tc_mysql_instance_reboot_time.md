Use this data source to query detailed information of mysql instance_reboot_time

Example Usage

```hcl
data "tencentcloud_mysql_instance_reboot_time" "instance_reboot_time" {
  instance_ids = ["cdb-fitq5t9h"]
}
```