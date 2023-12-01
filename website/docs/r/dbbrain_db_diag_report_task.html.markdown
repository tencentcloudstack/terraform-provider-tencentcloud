---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_db_diag_report_task"
sidebar_current: "docs-tencentcloud-resource-dbbrain_db_diag_report_task"
description: |-
  Provides a resource to create a dbbrain db_diag_report_task
---

# tencentcloud_dbbrain_db_diag_report_task

Provides a resource to create a dbbrain db_diag_report_task

## Example Usage

```hcl
resource "tencentcloud_dbbrain_db_diag_report_task" "db_diag_report_task" {
  instance_id    = "%s"
  start_time     = "%s"
  end_time       = "%s"
  send_mail_flag = 0
  product        = "mysql"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String, ForceNew) End time, such as 2020-11-09T14:00:00+08:00.
* `instance_id` - (Required, String, ForceNew) instance id.
* `product` - (Required, String, ForceNew) Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database CynosDB for MySQL.
* `send_mail_flag` - (Required, Int, ForceNew) Whether to send mail: 0 - no, 1 - yes.
* `start_time` - (Required, String, ForceNew) Start time, such as 2020-11-08T14:00:00+08:00.
* `contact_group` - (Optional, Set: [`Int`], ForceNew) An array of contact group IDs to receive mail from.
* `contact_person` - (Optional, Set: [`Int`], ForceNew) An array of contact IDs to receive emails from.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



