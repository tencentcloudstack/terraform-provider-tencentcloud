Use this data source to query detailed information of cls machine_group_configs

Example Usage

```hcl
resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-describe-mg-config-test"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type   = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

data "tencentcloud_cls_machine_group_configs" "machine_group_configs" {
  group_id = tencentcloud_cls_machine_group.group.id
}
```