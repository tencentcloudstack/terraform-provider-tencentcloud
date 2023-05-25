---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_machines"
sidebar_current: "docs-tencentcloud-datasource-cls_machines"
description: |-
  Use this data source to query detailed information of cls machines
---

# tencentcloud_cls_machines

Use this data source to query detailed information of cls machines

## Example Usage

```hcl
resource "tencentcloud_cls_machine_group" "group" {
  group_name        = "tf-describe-mg-test"
  service_logging   = true
  auto_update       = true
  update_end_time   = "19:05:00"
  update_start_time = "17:05:00"

  machine_group_type {
    type = "ip"
    values = [
      "192.168.1.1",
      "192.168.1.2",
    ]
  }
}

data "tencentcloud_cls_machines" "machines" {
  group_id = tencentcloud_cls_machine_group.group.id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Group id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `machines` - Info of Machines.
  * `auto_update` - if open auto update flag.
  * `err_code` - code of update operation.
  * `err_msg` - msg of update operation.
  * `ip` - ip of machine.
  * `offline_time` - offline time of machine.
  * `status` - status of machine.
  * `update_status` - machine update status.
  * `version` - current machine version.


